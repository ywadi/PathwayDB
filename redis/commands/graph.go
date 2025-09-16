package commands

import (
	"fmt"

	"github.com/ywadi/PathwayDB/models"
	"github.com/ywadi/PathwayDB/redis/protocol"
	"github.com/ywadi/PathwayDB/storage"
)

// GraphCommands handles graph-related Redis commands
type GraphCommands struct {
	storage storage.StorageEngine
}

// NewGraphCommands creates a new graph commands handler
func NewGraphCommands(storageEngine storage.StorageEngine) *GraphCommands {
	return &GraphCommands{
		storage: storageEngine,
	}
}

// Handle routes graph commands to their respective handlers
func (g *GraphCommands) Handle(command string, args []string) (*protocol.Response, error) {
	switch command {
	case "CREATE":
		return g.handleCreate(args)
	case "DELETE":
		return g.handleDelete(args)
	case "LIST":
		return g.handleList(args)
	case "GET":
		return g.handleGet(args)
	case "EXISTS":
		return g.handleExists(args)
	default:
		return nil, fmt.Errorf("unknown GRAPH command: %s", command)
	}
}

// handleCreate handles GRAPH.CREATE <name> [description]
func (g *GraphCommands) handleCreate(args []string) (*protocol.Response, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("GRAPH.CREATE requires at least 1 argument: name")
	}

	name := args[0]
	description := ""
	if len(args) > 1 {
		description = args[1]
	}

	graph := &models.Graph{
		ID:          models.GraphID(name),
		Name:        name,
		Description: description,
	}

	err := g.storage.CreateGraph(graph)
	if err != nil {
		return nil, fmt.Errorf("failed to create graph: %v", err)
	}

	return protocol.OK(), nil
}

// handleDelete handles GRAPH.DELETE <name>
func (g *GraphCommands) handleDelete(args []string) (*protocol.Response, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("GRAPH.DELETE requires exactly 1 argument: name")
	}

	name := args[0]
	err := g.storage.DeleteGraph(models.GraphID(name))
	if err != nil {
		return nil, fmt.Errorf("failed to delete graph: %v", err)
	}

	return protocol.OK(), nil
}

// handleList handles GRAPH.LIST
func (g *GraphCommands) handleList(args []string) (*protocol.Response, error) {
	graphs, err := g.storage.ListGraphs()
	if err != nil {
		return nil, fmt.Errorf("failed to list graphs: %v", err)
	}

	result := make([]string, 0, len(graphs)*2)
	for _, graph := range graphs {
		result = append(result, string(graph.ID), graph.Description)
	}

	return protocol.NewArrayResponse(result), nil
}

// handleGet handles GRAPH.GET <name>
func (g *GraphCommands) handleGet(args []string) (*protocol.Response, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("GRAPH.GET requires exactly 1 argument: name")
	}

	name := args[0]
	graph, err := g.storage.GetGraph(models.GraphID(name))
	if err != nil {
		return nil, fmt.Errorf("failed to get graph: %v", err)
	}

	if graph == nil {
		return protocol.NewNullResponse(), nil
	}

	// Return graph info as array: [id, name, description, node_count, edge_count]
	nodeCount, err := g.storage.CountNodes(graph.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to count nodes: %v", err)
	}

	edgeCount, err := g.storage.CountEdges(graph.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to count edges: %v", err)
	}
	
	result := []string{
		string(graph.ID),
		graph.Name,
		graph.Description,
		fmt.Sprintf("%d", nodeCount),
		fmt.Sprintf("%d", edgeCount),
	}

	return protocol.NewArrayResponse(result), nil
}

// handleExists handles GRAPH.EXISTS <name>
func (g *GraphCommands) handleExists(args []string) (*protocol.Response, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("GRAPH.EXISTS requires exactly 1 argument: name")
	}

	name := args[0]
	graph, err := g.storage.GetGraph(models.GraphID(name))
	if err != nil {
		return nil, fmt.Errorf("failed to check graph existence: %v", err)
	}

	if graph != nil {
		return protocol.NewIntResponse(1), nil
	}
	return protocol.NewIntResponse(0), nil
}
