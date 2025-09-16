package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/ywadi/PathwayDB/models"
)

// Key prefixes for different data types
const (
	GraphPrefix     = "g:"
	NodePrefix      = "n:"
	EdgePrefix      = "e:"
	NodeIndexPrefix = "ni:"
	EdgeIndexPrefix   = "ei:"
	TypeIndexPrefix   = "ti:"
	ExpiryIndexPrefix = "xi:"
)

// EncodeGraphKey creates a key for storing graph metadata
func EncodeGraphKey(graphID models.GraphID) []byte {
	return []byte(GraphPrefix + string(graphID))
}

// EncodeNodeKey creates a key for storing a node
func EncodeNodeKey(graphID models.GraphID, nodeID models.NodeID) []byte {
	return []byte(fmt.Sprintf("%s%s:%s", NodePrefix, graphID, nodeID))
}

// EncodeEdgeKey creates a key for storing an edge
func EncodeEdgeKey(graphID models.GraphID, edgeID models.EdgeID) []byte {
	return []byte(fmt.Sprintf("%s%s:%s", EdgePrefix, graphID, edgeID))
}

// EncodeNodeTypeIndexKey creates a key for indexing nodes by type
func EncodeNodeTypeIndexKey(graphID models.GraphID, nodeType models.NodeType, nodeID models.NodeID) []byte {
	return []byte(fmt.Sprintf("%sn:%s:%s:%s", TypeIndexPrefix, graphID, nodeType, nodeID))
}

// EncodeEdgeTypeIndexKey creates a key for indexing edges by type
func EncodeEdgeTypeIndexKey(graphID models.GraphID, edgeType models.EdgeType, edgeID models.EdgeID) []byte {
	return []byte(fmt.Sprintf("%se:%s:%s:%s", TypeIndexPrefix, graphID, edgeType, edgeID))
}

// EncodeNodeOutEdgeIndexKey creates a key for indexing outgoing edges from a node
func EncodeNodeOutEdgeIndexKey(graphID models.GraphID, nodeID models.NodeID, edgeID models.EdgeID) []byte {
	return []byte(fmt.Sprintf("%sout:%s:%s:%s", NodeIndexPrefix, graphID, nodeID, edgeID))
}

// EncodeNodeInEdgeIndexKey creates a key for indexing incoming edges to a node
func EncodeNodeInEdgeIndexKey(graphID models.GraphID, nodeID models.NodeID, edgeID models.EdgeID) []byte {
	return []byte(fmt.Sprintf("%sin:%s:%s:%s", NodeIndexPrefix, graphID, nodeID, edgeID))
}

// EncodeAttributeIndexKey creates a key for indexing nodes/edges by attribute
func EncodeAttributeIndexKey(graphID models.GraphID, entityType string, attrKey string, attrValue string, entityID string) []byte {
	return []byte(fmt.Sprintf("ai:%s:%s:%s:%s:%s", graphID, entityType, attrKey, attrValue, entityID))
}

// DecodeGraphID extracts graph ID from a graph key
func DecodeGraphID(key []byte) models.GraphID {
	keyStr := string(key)
	if strings.HasPrefix(keyStr, GraphPrefix) {
		return models.GraphID(keyStr[len(GraphPrefix):])
	}
	return ""
}

// DecodeNodeKey extracts graph ID and node ID from a node key
func DecodeNodeKey(key []byte) (models.GraphID, models.NodeID) {
	keyStr := string(key)
	if strings.HasPrefix(keyStr, NodePrefix) {
		parts := strings.Split(keyStr[len(NodePrefix):], ":")
		if len(parts) == 2 {
			return models.GraphID(parts[0]), models.NodeID(parts[1])
		}
	}
	return "", ""
}

// DecodeEdgeKey extracts graph ID and edge ID from an edge key
func DecodeEdgeKey(key []byte) (models.GraphID, models.EdgeID) {
	keyStr := string(key)
	if strings.HasPrefix(keyStr, EdgePrefix) {
		parts := strings.Split(keyStr[len(EdgePrefix):], ":")
		if len(parts) == 2 {
			return models.GraphID(parts[0]), models.EdgeID(parts[1])
		}
	}
	return "", ""
}

// CreateNodeIteratorPrefix creates a prefix for iterating over nodes in a graph
func CreateNodeIteratorPrefix(graphID models.GraphID) []byte {
	return []byte(fmt.Sprintf("%s%s:", NodePrefix, graphID))
}

// CreateEdgeIteratorPrefix creates a prefix for iterating over edges in a graph
func CreateEdgeIteratorPrefix(graphID models.GraphID) []byte {
	return []byte(fmt.Sprintf("%s%s:", EdgePrefix, graphID))
}

// CreateTypeIteratorPrefix creates a prefix for iterating over entities by type
func CreateTypeIteratorPrefix(graphID models.GraphID, entityType string, typeValue string) []byte {
	return []byte(fmt.Sprintf("%s%s:%s:%s:", TypeIndexPrefix, entityType, graphID, typeValue))
}

// EncodeExpiryIndexKey creates a key for the node expiration index.
func EncodeExpiryIndexKey(graphID models.GraphID, nodeID models.NodeID, expiresAt time.Time) []byte {
	// Use RFC3339 in UTC for lexicographical sorting of timestamps
	return []byte(fmt.Sprintf("%s%s:%s:%s", ExpiryIndexPrefix, expiresAt.UTC().Format(time.RFC3339), graphID, nodeID))
}

// DecodeExpiryIndexKey decodes the graph ID and node ID from an expiry index key.
func DecodeExpiryIndexKey(key []byte) (graphID models.GraphID, nodeID models.NodeID) {
	keyStr := strings.TrimPrefix(string(key), ExpiryIndexPrefix)

	// Find the last colon, which separates the node ID.
	lastColon := strings.LastIndex(keyStr, ":")
	if lastColon == -1 {
		return "", ""
	}
	nodeID = models.NodeID(keyStr[lastColon+1:])

	// Find the second to last colon, which separates the graph ID.
	remaining := keyStr[:lastColon]
	secondLastColon := strings.LastIndex(remaining, ":")
	if secondLastColon == -1 {
		return "", ""
	}
	graphID = models.GraphID(remaining[secondLastColon+1:])

	return graphID, nodeID
}

// CreateExpiryIteratorPrefix creates a prefix for iterating over the expiry index.
func CreateExpiryIteratorPrefix() []byte {
	return []byte(ExpiryIndexPrefix)
}
