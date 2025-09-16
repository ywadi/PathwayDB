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

// NodeCommands handles node-related Redis commands
type NodeCommands struct {
	storage storage.StorageEngine
}

// NewNodeCommands creates a new node commands handler
func NewNodeCommands(storageEngine storage.StorageEngine) *NodeCommands {
	return &NodeCommands{
		storage: storageEngine,
	}
}

// Handle routes node commands to their respective handlers
func (n *NodeCommands) Handle(command string, args []string) (*protocol.Response, error) {
	switch command {
	case "CREATE":
		return n.handleCreate(args)
	case "GET":
		return n.handleGet(args)
	case "UPDATE":
		return n.handleUpdate(args)
	case "DELETE":
		return n.handleDelete(args)
	case "FILTER":
		return n.handleFilter(args)
	case "LIST":
		return n.handleList(args)
	case "EXISTS":
		return n.handleExists(args)
	default:
		return nil, fmt.Errorf("unknown NODE command: %s", command)
	}
}

// handleCreate handles NODE.CREATE <graph> <id> <type> [attributes_json] [TTL <seconds>]
func (n *NodeCommands) handleCreate(args []string) (*protocol.Response, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("NODE.CREATE requires at least 3 arguments: graph, id, type")
	}

	graphID := args[0]
	nodeID := args[1]
	nodeType := args[2]

	attributes := make(map[string]interface{})
	var ttlSeconds int64 = -1

	// Parse optional arguments
	i := 3
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

	node := &models.Node{
		ID:         models.NodeID(nodeID),
		Type:       models.NodeType(nodeType),
		Attributes: attributes,
	}

	if ttlSeconds > 0 {
		expiresAt := time.Now().Add(time.Duration(ttlSeconds) * time.Second)
		node.ExpiresAt = &expiresAt
	}

	err := n.storage.CreateNode(models.GraphID(graphID), node)
	if err != nil {
		return nil, fmt.Errorf("failed to create node: %v", err)
	}

	return protocol.OK(), nil
}

// handleGet handles NODE.GET <graph> <id>
func (n *NodeCommands) handleGet(args []string) (*protocol.Response, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("NODE.GET requires exactly 2 arguments: graph, id")
	}

	graphID := args[0]
	nodeID := args[1]

	node, err := n.storage.GetNode(models.GraphID(graphID), models.NodeID(nodeID))
	if err != nil {
		return nil, fmt.Errorf("failed to get node: %v", err)
	}

	if node == nil {
		return protocol.NewNullResponse(), nil
	}

	// Serialize attributes
	attributesJSON, err := json.Marshal(node.Attributes)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize node attributes: %v", err)
	}

	expiresAtStr := ""
	if node.ExpiresAt != nil {
		expiresAtStr = node.ExpiresAt.Format(time.RFC3339)
	}

	result := []string{
		string(node.ID),
		string(node.Type),
		string(attributesJSON),
		expiresAtStr,
	}

	return protocol.NewArrayResponse(result), nil
}

// handleUpdate handles NODE.UPDATE <graph> <id> <attributes_json> [TTL <seconds>]
func (n *NodeCommands) handleUpdate(args []string) (*protocol.Response, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("NODE.UPDATE requires at least 3 arguments: graph, id, attributes_json")
	}

	graphID := args[0]
	nodeID := args[1]

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

	// Get the existing node first
	existingNode, err := n.storage.GetNode(models.GraphID(graphID), models.NodeID(nodeID))
	if err != nil {
		return nil, fmt.Errorf("failed to get node for update: %v", err)
	}
	if existingNode == nil {
		return nil, fmt.Errorf("node not found")
	}

	// Update attributes
	existingNode.Attributes = attributes

	if ttlSeconds >= 0 {
		if ttlSeconds == 0 {
			// TTL of 0 means remove expiration
			existingNode.ExpiresAt = nil
		} else {
			expiresAt := time.Now().Add(time.Duration(ttlSeconds) * time.Second)
			existingNode.ExpiresAt = &expiresAt
		}
	}

	err = n.storage.UpdateNode(models.GraphID(graphID), existingNode)
	if err != nil {
		return nil, fmt.Errorf("failed to update node: %v", err)
	}

	return protocol.OK(), nil
}

// handleDelete handles NODE.DELETE <graph> <id>
func (n *NodeCommands) handleDelete(args []string) (*protocol.Response, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("NODE.DELETE requires exactly 2 arguments: graph, id")
	}

	graphID := args[0]
	nodeID := args[1]

	err := n.storage.DeleteNode(models.GraphID(graphID), models.NodeID(nodeID))
	if err != nil {
		return nil, fmt.Errorf("failed to delete node: %v", err)
	}

	return protocol.OK(), nil
}

// handleFilter handles NODE.FILTER <graph> <attribute_key> <attribute_value>
func (n *NodeCommands) handleFilter(args []string) (*protocol.Response, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("NODE.FILTER requires exactly 3 arguments: graph, attribute_key, attribute_value")
	}

	graphID := args[0]
	attrKey := args[1]
	attrValue := args[2]

	// Attempt to unmarshal the value as JSON, if it fails, use it as a string
	var value interface{}
	if err := json.Unmarshal([]byte(attrValue), &value); err != nil {
		value = attrValue
	}

	nodes, err := n.storage.FindNodesByAttribute(models.GraphID(graphID), attrKey, value)
	if err != nil {
		return nil, fmt.Errorf("failed to filter nodes by attribute: %v", err)
	}

	// Format response as array of node data
	result := make([]string, 0, len(nodes)*3)
	for _, node := range nodes {
		attributesJSON, err := json.Marshal(node.Attributes)
		if err != nil {
			// Log or handle this error, maybe skip the node
			continue
		}
		result = append(result, string(node.ID), string(node.Type), string(attributesJSON))
	}

	return protocol.NewArrayResponse(result), nil
}

// handleList handles NODE.LIST <graph>
func (n *NodeCommands) handleList(args []string) (*protocol.Response, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("NODE.LIST requires exactly 1 argument: graph")
	}

	graphID := args[0]
	// Get all nodes in the graph using ListNodes instead
	nodes, err := n.storage.ListNodes(models.GraphID(graphID))
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %v", err)
	}

	if nodes == nil {
		return protocol.NewArrayResponse([]string{}), nil
	}

	// Return node IDs and types
	result := make([]string, 0, len(nodes)*2)
	for _, node := range nodes {
		result = append(result, string(node.ID), string(node.Type))
	}

	return protocol.NewArrayResponse(result), nil
}

// handleExists handles NODE.EXISTS <graph> <id>
func (n *NodeCommands) handleExists(args []string) (*protocol.Response, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("NODE.EXISTS requires exactly 2 arguments: graph, id")
	}

	graphID := args[0]
	nodeID := args[1]

	node, err := n.storage.GetNode(models.GraphID(graphID), models.NodeID(nodeID))
	if err != nil {
		return nil, fmt.Errorf("failed to check node existence: %v", err)
	}

	if node != nil {
		return protocol.NewIntResponse(1), nil
	}
	return protocol.NewIntResponse(0), nil
}
