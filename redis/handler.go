package redis

import (
	"fmt"
	"strings"

	"github.com/ywadi/PathwayDB/redis/commands"
	"github.com/ywadi/PathwayDB/redis/protocol"
	"github.com/ywadi/PathwayDB/storage"
)

// Response is an alias for protocol.Response for convenience
type Response = protocol.Response

// CommandHandler handles Redis command routing and execution
type CommandHandler struct {
	storage      storage.StorageEngine
	graphCmd     *commands.GraphCommands
	nodeCmd      *commands.NodeCommands
	edgeCmd      *commands.EdgeCommands
	analysisCmd  *commands.AnalysisCommands
}

// NewCommandHandler creates a new command handler
func NewCommandHandler(storageEngine storage.StorageEngine) *CommandHandler {
	return &CommandHandler{
		storage:     storageEngine,
		graphCmd:    commands.NewGraphCommands(storageEngine),
		nodeCmd:     commands.NewNodeCommands(storageEngine),
		edgeCmd:     commands.NewEdgeCommands(storageEngine),
		analysisCmd: commands.NewAnalysisCommands(storageEngine),
	}
}

// Handle routes and executes Redis commands
func (h *CommandHandler) Handle(command string, args []string) (*Response, error) {
	// Split command by dots for namespaced commands (e.g., GRAPH.CREATE)
	parts := strings.Split(command, ".")
	
	switch parts[0] {
	case "PING":
		return h.handlePing(args)
	case "INFO":
		return h.handleInfo(args)
	case "GRAPH":
		if len(parts) < 2 {
			return nil, fmt.Errorf("incomplete GRAPH command")
		}
		return h.graphCmd.Handle(parts[1], args)
	case "NODE":
		if len(parts) < 2 {
			return nil, fmt.Errorf("incomplete NODE command")
		}
		return h.nodeCmd.Handle(parts[1], args)
	case "EDGE":
		if len(parts) < 2 {
			return nil, fmt.Errorf("incomplete EDGE command")
		}
		return h.edgeCmd.Handle(parts[1], args)
	case "ANALYSIS":
		if len(parts) < 2 {
			return nil, fmt.Errorf("incomplete ANALYSIS command")
		}
		return h.analysisCmd.Handle(parts[1], args)
	default:
		return nil, fmt.Errorf("unknown command: %s", command)
	}
}

// handlePing handles the PING command
func (h *CommandHandler) handlePing(args []string) (*Response, error) {
	if len(args) == 0 {
		return protocol.NewStringResponse("PONG"), nil
	}
	return protocol.NewBulkResponse(args[0]), nil
}

// handleInfo handles the INFO command
func (h *CommandHandler) handleInfo(args []string) (*Response, error) {
	info := []string{
		"# PathwayDB",
		"version:1.0.0",
		"redis_protocol:enabled",
		"storage_engine:badger",
	}
	
	return protocol.NewBulkResponse(strings.Join(info, "\r\n")), nil
}
