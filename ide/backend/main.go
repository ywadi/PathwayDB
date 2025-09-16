package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

type WebSocketMessage struct {
	ID        string   `json:"id"`
	Command   string   `json:"command"`
	Args      []string `json:"args"`
	Timestamp int64    `json:"timestamp"`
}

type WebSocketResponse struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
	Timestamp int64       `json:"timestamp"`
}

type ConnectionPool struct {
	redisAddr string
	pool      chan net.Conn
	maxSize   int
	mutex     sync.Mutex
	closed    bool
}

func NewConnectionPool(redisAddr string, maxSize int) *ConnectionPool {
	return &ConnectionPool{
		redisAddr: redisAddr,
		pool:      make(chan net.Conn, maxSize),
		maxSize:   maxSize,
	}
}

func (cp *ConnectionPool) GetConnection() (net.Conn, error) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	
	if cp.closed {
		return nil, fmt.Errorf("connection pool is closed")
	}
	
	select {
	case conn := <-cp.pool:
		// Test if connection is still alive
		conn.SetReadDeadline(time.Now().Add(1 * time.Millisecond))
		buf := make([]byte, 1)
		_, err := conn.Read(buf)
		conn.SetReadDeadline(time.Time{}) // Clear deadline
		
		if err == nil {
			// Connection has data, put it back and create new one
			return net.Dial("tcp", cp.redisAddr)
		} else if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			// Connection is alive (timeout as expected)
			return conn, nil
		} else {
			// Connection is dead, create new one
			conn.Close()
			return net.Dial("tcp", cp.redisAddr)
		}
	default:
		// No connection available, create new one
		return net.Dial("tcp", cp.redisAddr)
	}
}

func (cp *ConnectionPool) ReturnConnection(conn net.Conn) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	
	if cp.closed {
		conn.Close()
		return
	}
	
	select {
	case cp.pool <- conn:
		// Successfully returned to pool
	default:
		// Pool is full, close the connection
		conn.Close()
	}
}

func (cp *ConnectionPool) Close() {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	
	cp.closed = true
	close(cp.pool)
	
	// Close all connections in pool
	for conn := range cp.pool {
		conn.Close()
	}
}

type RedisProxy struct {
	redisAddr string
	connPool  *ConnectionPool
}

func NewRedisProxy(redisAddr string) *RedisProxy {
	return &RedisProxy{
		redisAddr: redisAddr,
		connPool:  NewConnectionPool(redisAddr, 10), // Pool of 10 connections
	}
}

func (rp *RedisProxy) ExecuteCommand(command string, args []string) (*WebSocketResponse, error) {
	// Get connection from pool
	conn, err := rp.connPool.GetConnection()
	if err != nil {
		return nil, err
	}
	
	// Return connection to pool when done (or close if error)
	defer func() {
		if err != nil {
			conn.Close()
		} else {
			rp.connPool.ReturnConnection(conn)
		}
	}()

	// Build Redis command in RESP protocol
	cmdParts := append([]string{command}, args...)
	redisCmd := fmt.Sprintf("*%d\r\n", len(cmdParts))
	for _, part := range cmdParts {
		redisCmd += fmt.Sprintf("$%d\r\n%s\r\n", len(part), part)
	}

	// Send command
	_, err = conn.Write([]byte(redisCmd))
	if err != nil {
		return &WebSocketResponse{
			Type:      "error",
			Value:     err.Error(),
			Timestamp: time.Now().UnixMilli(),
		}, nil
	}

	// Read response
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return &WebSocketResponse{
			Type:      "error",
			Value:     err.Error(),
			Timestamp: time.Now().UnixMilli(),
		}, nil
	}

	responseStr := string(buf[:n])
	response := &WebSocketResponse{
		Timestamp: time.Now().UnixMilli(),
	}

	// Parse Redis response
	if len(responseStr) == 0 {
		response.Type = "null"
		response.Value = nil
	} else {
		switch responseStr[0] {
		case '+':
			response.Type = "string"
			response.Value = strings.TrimSpace(responseStr[1:])
		case '-':
			response.Type = "error"
			response.Value = strings.TrimSpace(responseStr[1:])
		case ':':
			response.Type = "int"
			if val, err := strconv.ParseInt(strings.TrimSpace(responseStr[1:]), 10, 64); err == nil {
				response.Value = val
			} else {
				response.Value = 0
			}
		case '$':
			response.Type = "bulk"
			lines := strings.Split(responseStr, "\r\n")
			if len(lines) > 1 {
				response.Value = lines[1]
			} else {
				response.Value = ""
			}
		case '*':
			response.Type = "array"
			lines := strings.Split(responseStr, "\r\n")
			var result []string
			for i := 1; i < len(lines); i += 2 {
				if i+1 < len(lines) && lines[i][0] == '$' {
					result = append(result, lines[i+1])
				}
			}
			response.Value = result
		default:
			response.Type = "string"
			response.Value = strings.TrimSpace(responseStr)
		}
	}

	return response, nil
}

func (rp *RedisProxy) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("WebSocket client connected: %s", conn.RemoteAddr())

	for {
		var msg WebSocketMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		log.Printf("Received command: %s %v", msg.Command, msg.Args)

		// Execute Redis command
		response, err := rp.ExecuteCommand(msg.Command, msg.Args)
		if err != nil {
			response = &WebSocketResponse{
				Type:      "error",
				Value:     err.Error(),
				Timestamp: time.Now().UnixMilli(),
			}
		}

		// Set response ID to match request
		response.ID = msg.ID

		// Send response back to client
		if err := conn.WriteJSON(response); err != nil {
			log.Printf("Failed to send response: %v", err)
			break
		}
	}

	log.Printf("WebSocket client disconnected: %s", conn.RemoteAddr())
}

func (rp *RedisProxy) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"redis":  rp.redisAddr,
		"time":   time.Now().Unix(),
	})
}

func (rp *RedisProxy) handleListDocs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get the project root directory
	execPath, err := os.Executable()
	if err != nil {
		http.Error(w, "Failed to determine executable path", http.StatusInternalServerError)
		return
	}
	
	// Navigate to project root from wherever the executable is
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(execPath)))
	// If running with go run, use current working directory approach
	if strings.Contains(execPath, "go-build") {
		wd, err := os.Getwd()
		if err != nil {
			http.Error(w, "Failed to get working directory", http.StatusInternalServerError)
			return
		}
		// If we're in the backend directory, go up two levels
		if strings.HasSuffix(wd, "ide/backend") {
			projectRoot = filepath.Dir(filepath.Dir(wd))
		} else {
			projectRoot = wd
		}
	}
	
	docsDir := filepath.Join(projectRoot, "docs")
	files, err := ioutil.ReadDir(docsDir)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read docs directory: %s", docsDir), http.StatusInternalServerError)
		return
	}

	var docFiles []string
	// Add root README
	docFiles = append(docFiles, "README.md")

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			docFiles = append(docFiles, file.Name())
		}
	}

	json.NewEncoder(w).Encode(docFiles)
}

func (rp *RedisProxy) handleGetDoc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	filename := strings.TrimPrefix(r.URL.Path, "/api/docs/")
	if filename == "" {
		http.Error(w, "Filename not specified", http.StatusBadRequest)
		return
	}

	// Get the project root directory
	execPath, err := os.Executable()
	if err != nil {
		http.Error(w, "Failed to determine executable path", http.StatusInternalServerError)
		return
	}
	
	// Navigate to project root from wherever the executable is
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(execPath)))
	// If running with go run, use current working directory approach
	if strings.Contains(execPath, "go-build") {
		wd, err := os.Getwd()
		if err != nil {
			http.Error(w, "Failed to get working directory", http.StatusInternalServerError)
			return
		}
		// If we're in the backend directory, go up two levels
		if strings.HasSuffix(wd, "ide/backend") {
			projectRoot = filepath.Dir(filepath.Dir(wd))
		} else {
			projectRoot = wd
		}
	}

	// Determine file path
	var filePath string
	if filename == "README.md" {
		filePath = filepath.Join(projectRoot, "README.md")
	} else {
		filePath = filepath.Join(projectRoot, "docs", filename)
	}

	// Security: ensure the path is clean and within the project
	cleanPath, err := filepath.Abs(filePath)
	if err != nil {
		http.Error(w, "Invalid file path", http.StatusInternalServerError)
		return
	}
	// A basic check to prevent directory traversal
	if !strings.Contains(cleanPath, "PathwayDB") {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	md, err := ioutil.ReadFile(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
		} else {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
		}
		return
	}

	// Convert markdown to HTML
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)
	html := markdown.ToHTML(md, parser, nil)

	w.Write(html)
}

// getEnv reads an environment variable or returns a fallback value.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	// Configuration with environment variable overrides
	websocketAddr := getEnv("WEBSOCKET_ADDR", ":8081")
	redisAddrEnv := getEnv("REDIS_ADDR", "localhost:6379")

	var (
		addr      = flag.String("addr", websocketAddr, "WebSocket server address")
		redisAddr = flag.String("redis", redisAddrEnv, "Redis server address")
	)
	flag.Parse()

	proxy := NewRedisProxy(*redisAddr)
	
	// Cleanup connection pool on shutdown
	defer proxy.connPool.Close()

	// WebSocket endpoint
	http.HandleFunc("/ws", proxy.handleWebSocket)

	// Health check endpoint
	http.HandleFunc("/health", proxy.handleHealth)

	// Documentation endpoint
	http.HandleFunc("/api/docs/", proxy.handleGetDoc)
	http.HandleFunc("/api/docs", proxy.handleListDocs)

	// Serve static files (for development)
	http.Handle("/", http.FileServer(http.Dir("../frontend/build/")))

	log.Printf("PathwayDB IDE WebSocket server starting on %s", *addr)
	log.Printf("Connecting to Redis server at %s with connection pool (max 10 connections)", *redisAddr)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
