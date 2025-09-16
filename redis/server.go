package redis

import (
	"log"
	"strings"
	"sync"

	"github.com/tidwall/redcon"
	"github.com/ywadi/PathwayDB/redis/protocol"
	"github.com/ywadi/PathwayDB/storage"
)

// Server represents the Redis protocol server for PathwayDB
type Server struct {
	config  *Config
	storage storage.StorageEngine
	handler *CommandHandler
	mu      sync.RWMutex
	running bool
}

// NewServer creates a new Redis protocol server
func NewServer(config *Config, storageEngine storage.StorageEngine) *Server {
	server := &Server{
		config:  config,
		storage: storageEngine,
		handler: NewCommandHandler(storageEngine),
	}
	return server
}

// Start starts the Redis protocol server
func (s *Server) Start() error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return nil
	}
	s.running = true
	s.mu.Unlock()

	log.Printf("Starting PathwayDB Redis server on %s", s.config.Address)

	return redcon.ListenAndServe(s.config.Address,
		s.handleConnection,
		s.handleAccept,
		s.handleClosed,
	)
}

// Stop stops the Redis protocol server
func (s *Server) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = false
}

// handleConnection handles incoming Redis commands
func (s *Server) handleConnection(conn redcon.Conn, cmd redcon.Command) {
	// Parse command
	if len(cmd.Args) == 0 {
		conn.WriteError("ERR empty command")
		return
	}

	command := strings.ToUpper(string(cmd.Args[0]))
	args := make([]string, len(cmd.Args)-1)
	for i, arg := range cmd.Args[1:] {
		args[i] = string(arg)
	}

	// Route command to handler
	response, err := s.handler.Handle(command, args)
	if err != nil {
		conn.WriteError("ERR " + err.Error())
		return
	}

	// Write response
	s.writeResponse(conn, response)
}

// handleAccept handles new client connections
func (s *Server) handleAccept(conn redcon.Conn) bool {
	log.Printf("Client connected: %s", conn.RemoteAddr())
	return true
}

// handleClosed handles client disconnections
func (s *Server) handleClosed(conn redcon.Conn, err error) {
	if err != nil {
		log.Printf("Client disconnected with error: %s, error: %v", conn.RemoteAddr(), err)
	} else {
		log.Printf("Client disconnected: %s", conn.RemoteAddr())
	}
}

// writeResponse writes a response to the Redis connection
func (s *Server) writeResponse(conn redcon.Conn, response *Response) {
	switch response.Type {
	case protocol.ResponseTypeString:
		conn.WriteString(response.StringValue)
	case protocol.ResponseTypeInt:
		conn.WriteInt64(response.IntValue)
	case protocol.ResponseTypeArray:
		conn.WriteArray(len(response.ArrayValue))
		for _, item := range response.ArrayValue {
			conn.WriteBulkString(item)
		}
	case protocol.ResponseTypeNestedArray:
		conn.WriteArray(len(response.NestedArrayValue))
		for _, subArray := range response.NestedArrayValue {
			if sa, ok := subArray.([]string); ok {
				conn.WriteArray(len(sa))
				for _, item := range sa {
					conn.WriteBulkString(item)
				}
			} else {
				conn.WriteError("ERR invalid nested array format")
			}
		}
	case protocol.ResponseTypeBulk:
		conn.WriteBulkString(response.StringValue)
	case protocol.ResponseTypeNull:
		conn.WriteNull()
	case protocol.ResponseTypeError:
		conn.WriteError(response.StringValue)
	default:
		conn.WriteError("ERR unknown response type")
	}
}

// IsRunning returns whether the server is currently running
func (s *Server) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}
