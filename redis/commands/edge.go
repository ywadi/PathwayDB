package commands

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ywadi/PathwayDB/models"
	"github.com/ywadi/PathwayDB/redis/protocol"
	"github.com/ywadi/PathwayDB/storage"
)

// EdgeCommands handles edge-related Redis commands
type EdgeCommands struct {
	storage storage.StorageEngine
}

// NewEdgeCommands creates a new edge commands handler
func NewEdgeCommands(storageEngine storage.StorageEngine) *EdgeCommands {
	return &EdgeCommands{
		storage: storageEngine,
	}
}

// Handle routes edge commands to their respective handlers
func (e *EdgeCommands) Handle(command string, args []string) (*protocol.Response, error) {
	switch command {
	case "CREATE":
		return e.handleCreate(args)
	case "GET":
		return e.handleGet(args)
	case "UPDATE":
		return e.handleUpdate(args)
	case "DELETE":
		return e.handleDelete(args)
	case "FILTER":
		return e.handleFilter(args)
	case "NEIGHBORS":
		return e.handleNeighbors(args)
	case "LIST":
		return e.handleList(args)
	case "EXISTS":
		return e.handleExists(args)
	default:
		return nil, fmt.Errorf("unknown EDGE command: %s", command)
	}
}

// handleCreate handles EDGE.CREATE <graph> <id> <from> <to> <type> [attributes_json] [TTL <seconds>]
func (e *EdgeCommands) handleCreate(args []string) (*protocol.Response, error) {
	if len(args) < 5 {
		return nil, fmt.Errorf("EDGE.CREATE requires at least 5 arguments: graph, id, from, to, type")
	}

	graphID := args[0]
	edgeID := args[1]
	fromNodeID := args[2]
	toNodeID := args[3]
	edgeType := args[4]

	attributes := make(map[string]interface{})
	var ttlSeconds int64 = -1

	// Parse optional arguments
	i := 5
	for i < len(args) {
		switch strings.ToUpper(args[i]) {
		case "TTL":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("TTL option requires a value")
			}
			ttl, err := strconv.ParseInt(args[i+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid TTL value: %v", err)
			}
			ttlSeconds = ttl
			i += 2
		default:
			// Assume it's the attributes JSON
			if err := json.Unmarshal([]byte(args[i]), &attributes); err != nil {
				return nil, fmt.Errorf("invalid attributes JSON: %v", err)
			}
			i++
		}
	}

	edge := &models.Edge{
		ID:         models.EdgeID(edgeID),
		FromNodeID: models.NodeID(fromNodeID),
		ToNodeID:   models.NodeID(toNodeID),
		Type:       models.EdgeType(edgeType),
		Attributes: attributes,
	}

	if ttlSeconds > 0 {
		expiresAt := time.Now().Add(time.Duration(ttlSeconds) * time.Second)
		edge.ExpiresAt = &expiresAt
	}

	err := e.storage.CreateEdge(models.GraphID(graphID), edge)
	if err != nil {
		return nil, fmt.Errorf("failed to create edge: %v", err)
	}

	return protocol.OK(), nil
}

// handleGet handles EDGE.GET <graph> <id>
func (e *EdgeCommands) handleGet(args []string) (*protocol.Response, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("EDGE.GET requires exactly 2 arguments: graph, id")
	}

	graphID := args[0]
	edgeID := args[1]

	edge, err := e.storage.GetEdge(models.GraphID(graphID), models.EdgeID(edgeID))
	if err != nil {
		return nil, fmt.Errorf("failed to get edge: %v", err)
	}

	if edge == nil {
		return protocol.NewNullResponse(), nil
	}

	// Serialize attributes
	attributesJSON, err := json.Marshal(edge.Attributes)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize edge attributes: %v", err)
	}

	expiresAtStr := ""
	if edge.ExpiresAt != nil {
		expiresAtStr = edge.ExpiresAt.Format(time.RFC3339)
	}

	result := []string{
		string(edge.ID),
		string(edge.FromNodeID),
		string(edge.ToNodeID),
		string(edge.Type),
		string(attributesJSON),
		expiresAtStr,
	}

	return protocol.NewArrayResponse(result), nil
}

// handleUpdate handles EDGE.UPDATE <graph> <id> <attributes_json> [TTL <seconds>]
func (e *EdgeCommands) handleUpdate(args []string) (*protocol.Response, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("EDGE.UPDATE requires at least 3 arguments: graph, id, attributes_json")
	}

	graphID := args[0]
	edgeID := args[1]

	// Parse new attributes
	var attributes map[string]interface{}
	if err := json.Unmarshal([]byte(args[2]), &attributes); err != nil {
		return nil, fmt.Errorf("invalid attributes JSON: %v", err)
	}

	var ttlSeconds int64 = -1
	// Parse optional TTL
	if len(args) > 4 && strings.ToUpper(args[3]) == "TTL" {
		ttl, err := strconv.ParseInt(args[4], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid TTL value: %v", err)
		}
		ttlSeconds = ttl
	}

	// Get the existing edge first
	existingEdge, err := e.storage.GetEdge(models.GraphID(graphID), models.EdgeID(edgeID))
	if err != nil {
		return nil, fmt.Errorf("failed to get edge for update: %v", err)
	}
	if existingEdge == nil {
		return nil, fmt.Errorf("edge not found")
	}

	// Update attributes
	existingEdge.Attributes = attributes

	if ttlSeconds >= 0 {
		if ttlSeconds == 0 {
			// TTL of 0 means remove expiration
			existingEdge.ExpiresAt = nil
		} else {
			expiresAt := time.Now().Add(time.Duration(ttlSeconds) * time.Second)
			existingEdge.ExpiresAt = &expiresAt
		}
	}

	err = e.storage.UpdateEdge(models.GraphID(graphID), existingEdge)
	if err != nil {
		return nil, fmt.Errorf("failed to update edge: %v", err)
	}

	return protocol.OK(), nil
}

// handleDelete handles EDGE.DELETE <graph> <id>
func (e *EdgeCommands) handleDelete(args []string) (*protocol.Response, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("EDGE.DELETE requires exactly 2 arguments: graph, id")
	}

	graphID := args[0]
	edgeID := args[1]

	err := e.storage.DeleteEdge(models.GraphID(graphID), models.EdgeID(edgeID))
	if err != nil {
		return nil, fmt.Errorf("failed to delete edge: %v", err)
	}

	return protocol.OK(), nil
}

// handleFilter handles EDGE.FILTER <graph> <attribute_key> <attribute_value>
func (e *EdgeCommands) handleFilter(args []string) (*protocol.Response, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("EDGE.FILTER requires exactly 3 arguments: graph, attribute_key, attribute_value")
	}

	graphID := args[0]
	attrKey := args[1]
	attrValue := args[2]

	// Attempt to unmarshal the value as JSON, if it fails, use it as a string
	var value interface{}
	if err := json.Unmarshal([]byte(attrValue), &value); err != nil {
		value = attrValue
	}

	edges, err := e.storage.FindEdgesByAttribute(models.GraphID(graphID), attrKey, value)
	if err != nil {
		return nil, fmt.Errorf("failed to filter edges by attribute: %v", err)
	}

	// Format response as array of edge data
	result := make([]string, 0, len(edges)*5)
	for _, edge := range edges {
		attributesJSON, err := json.Marshal(edge.Attributes)
		if err != nil {
			// Log or handle this error, maybe skip the edge
			continue
		}
		result = append(result, string(edge.ID), string(edge.FromNodeID), string(edge.ToNodeID), string(edge.Type), string(attributesJSON))
	}

	return protocol.NewArrayResponse(result), nil
}

// handleNeighbors handles EDGE.NEIGHBORS <graph> <node_id> [direction] [FORMAT simple|detailed]
// direction can be: "in", "out", "both" (default: "both")
func (e *EdgeCommands) handleNeighbors(args []string) (*protocol.Response, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("EDGE.NEIGHBORS requires at least 2 arguments: graph, node_id")
	}

	graphID := args[0]
	nodeID := args[1]
	direction := "both"
	format := "detailed" // Default to detailed format

	// Parse optional arguments
	for i := 2; i < len(args); i++ {
		if args[i] == "FORMAT" && i+1 < len(args) {
			i++
			if args[i] != "simple" && args[i] != "detailed" {
				return nil, fmt.Errorf("invalid FORMAT: %s (must be 'simple' or 'detailed')", args[i])
			}
			format = args[i]
		} else if args[i] == "in" || args[i] == "out" || args[i] == "both" {
			direction = args[i]
		} else if args[i] != "FORMAT" {
			return nil, fmt.Errorf("invalid argument: %s", args[i])
		}
	}

	type NeighborInfo struct {
		Node     *models.Node
		Edge     *models.Edge
		Direction string
	}

	var neighborInfos []NeighborInfo

	switch direction {
	case "in":
		incomingEdges, err := e.storage.GetIncomingEdges(models.GraphID(graphID), models.NodeID(nodeID))
		if err != nil {
			return nil, fmt.Errorf("failed to get incoming edges: %v", err)
		}
		for _, edge := range incomingEdges {
			node, err := e.storage.GetNode(models.GraphID(graphID), edge.FromNodeID)
			if err == nil && node != nil {
				neighborInfos = append(neighborInfos, NeighborInfo{
					Node:      node,
					Edge:      edge,
					Direction: "in",
				})
			}
		}
	case "out":
		outgoingEdges, err := e.storage.GetOutgoingEdges(models.GraphID(graphID), models.NodeID(nodeID))
		if err != nil {
			return nil, fmt.Errorf("failed to get outgoing edges: %v", err)
		}
		for _, edge := range outgoingEdges {
			node, err := e.storage.GetNode(models.GraphID(graphID), edge.ToNodeID)
			if err == nil && node != nil {
				neighborInfos = append(neighborInfos, NeighborInfo{
					Node:      node,
					Edge:      edge,
					Direction: "out",
				})
			}
		}
	case "both":
		// Get incoming edges
		incomingEdges, err := e.storage.GetIncomingEdges(models.GraphID(graphID), models.NodeID(nodeID))
		if err != nil {
			return nil, fmt.Errorf("failed to get incoming edges: %v", err)
		}
		for _, edge := range incomingEdges {
			node, err := e.storage.GetNode(models.GraphID(graphID), edge.FromNodeID)
			if err == nil && node != nil {
				neighborInfos = append(neighborInfos, NeighborInfo{
					Node:      node,
					Edge:      edge,
					Direction: "in",
				})
			}
		}
		
		// Get outgoing edges
		outgoingEdges, err := e.storage.GetOutgoingEdges(models.GraphID(graphID), models.NodeID(nodeID))
		if err != nil {
			return nil, fmt.Errorf("failed to get outgoing edges: %v", err)
		}
		for _, edge := range outgoingEdges {
			node, err := e.storage.GetNode(models.GraphID(graphID), edge.ToNodeID)
			if err == nil && node != nil {
				neighborInfos = append(neighborInfos, NeighborInfo{
					Node:      node,
					Edge:      edge,
					Direction: "out",
				})
			}
		}
	}

	// Simple format with nodeid:nodetype
	if format == "simple" {
		response := make([]string, len(neighborInfos))
		for i, info := range neighborInfos {
			response[i] = string(info.Node.ID) + ":" + string(info.Node.Type)
		}
		return protocol.NewArrayResponse(response), nil
	}

	// Enhanced detailed format with pipe-delimited format
	result := make([]string, 0, len(neighborInfos)+1)
	result = append(result, fmt.Sprintf("%d", len(neighborInfos)))
	
	for _, info := range neighborInfos {
		// Format: neighbor_node_id:neighbor_node_type->connecting_edge_id:connecting_edge_type->direction
		neighborStr := fmt.Sprintf("%s:%s->%s:%s->%s",
			string(info.Node.ID),
			string(info.Node.Type),
			string(info.Edge.ID),
			string(info.Edge.Type),
			info.Direction)
		result = append(result, neighborStr)
	}

	return protocol.NewArrayResponse(result), nil
}

// handleList handles EDGE.LIST <graph>
func (e *EdgeCommands) handleList(args []string) (*protocol.Response, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("EDGE.LIST requires exactly 1 argument: graph")
	}

	graphID := args[0]
	// Get all edges in the graph using ListEdges instead
	edges, err := e.storage.ListEdges(models.GraphID(graphID))
	if err != nil {
		return nil, fmt.Errorf("failed to get graph: %v", err)
	}

	if edges == nil {
		return protocol.NewArrayResponse([]string{}), nil
	}

	// Return edge IDs, from/to nodes, and types
	result := make([]string, 0, len(edges)*4)
	for _, edge := range edges {
		result = append(result, string(edge.ID), string(edge.FromNodeID), string(edge.ToNodeID), string(edge.Type))
	}

	return protocol.NewArrayResponse(result), nil
}

// handleExists handles EDGE.EXISTS <graph> <id>
func (e *EdgeCommands) handleExists(args []string) (*protocol.Response, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("EDGE.EXISTS requires exactly 2 arguments: graph, id")
	}

	graphID := args[0]
	edgeID := args[1]

	edge, err := e.storage.GetEdge(models.GraphID(graphID), models.EdgeID(edgeID))
	if err != nil {
		return nil, fmt.Errorf("failed to check edge existence: %v", err)
	}

	if edge != nil {
		return protocol.NewIntResponse(1), nil
	}
	return protocol.NewIntResponse(0), nil
}
