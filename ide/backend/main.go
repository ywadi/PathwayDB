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

type RedisProxy struct {
	redisAddr string
}

func NewRedisProxy(redisAddr string) *RedisProxy {
	return &RedisProxy{
		redisAddr: redisAddr,
	}
}

func (rp *RedisProxy) ExecuteCommand(command string, args []string) (*WebSocketResponse, error) {
	// Connect to Redis server via TCP
	conn, err := net.Dial("tcp", rp.redisAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

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

	docsDir := "../../docs"
	files, err := ioutil.ReadDir(docsDir)
	if err != nil {
		http.Error(w, "Failed to read docs directory", http.StatusInternalServerError)
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

	// Determine file path
	var filePath string
	if filename == "README.md" {
		filePath = "../../README.md"
	} else {
		filePath = filepath.Join("../../docs", filename)
	}

	// Security: ensure the path is clean and within the project
	cleanPath, err := filepath.Abs(filePath)
	if err != nil {
		http.Error(w, "Invalid file path", http.StatusInternalServerError)
		return
	}
	// A basic check to prevent directory traversal
	if !strings.Contains(cleanPath, "/Development/PathwayDB/") {
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
	log.Printf("Connecting to Redis server at %s", *redisAddr)

	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
