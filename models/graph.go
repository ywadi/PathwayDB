package models

import (
	"encoding/json"
	"time"
)

// NodeType represents the type of a node
type NodeType string

// EdgeType represents the type of an edge
type EdgeType string

// NodeID represents a unique identifier for a node
type NodeID string

// EdgeID represents a unique identifier for an edge
type EdgeID string

// GraphID represents a unique identifier for a graph
type GraphID string

// Attributes represents key-value pairs for node/edge attributes
type Attributes map[string]interface{}

// Node represents a vertex in the graph
type Node struct {
	ID         NodeID     `json:"id"`
	Type       NodeType   `json:"type"`
	Attributes Attributes `json:"attributes"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

// Edge represents a connection between two nodes
type Edge struct {
	ID         EdgeID     `json:"id"`
	Type       EdgeType   `json:"type"`
	FromNodeID NodeID     `json:"from_node_id"`
	ToNodeID   NodeID     `json:"to_node_id"`
	Attributes Attributes `json:"attributes"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

// Graph represents a collection of nodes and edges
type Graph struct {
	ID          GraphID   `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToJSON converts a node to JSON bytes
func (n *Node) ToJSON() ([]byte, error) {
	return json.Marshal(n)
}

// FromJSON populates a node from JSON bytes
func (n *Node) FromJSON(data []byte) error {
	return json.Unmarshal(data, n)
}

// ToJSON converts an edge to JSON bytes
func (e *Edge) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// FromJSON populates an edge from JSON bytes
func (e *Edge) FromJSON(data []byte) error {
	return json.Unmarshal(data, e)
}

// ToJSON converts a graph to JSON bytes
func (g *Graph) ToJSON() ([]byte, error) {
	return json.Marshal(g)
}

// FromJSON populates a graph from JSON bytes
func (g *Graph) FromJSON(data []byte) error {
	return json.Unmarshal(data, g)
}

// HasAttribute checks if a node has a specific attribute
func (n *Node) HasAttribute(key string) bool {
	_, exists := n.Attributes[key]
	return exists
}

// GetAttribute gets an attribute value from a node
func (n *Node) GetAttribute(key string) (interface{}, bool) {
	value, exists := n.Attributes[key]
	return value, exists
}

// SetAttribute sets an attribute on a node
func (n *Node) SetAttribute(key string, value interface{}) {
	if n.Attributes == nil {
		n.Attributes = make(Attributes)
	}
	n.Attributes[key] = value
	n.UpdatedAt = time.Now()
}

// HasAttribute checks if an edge has a specific attribute
func (e *Edge) HasAttribute(key string) bool {
	_, exists := e.Attributes[key]
	return exists
}

// GetAttribute gets an attribute value from an edge
func (e *Edge) GetAttribute(key string) (interface{}, bool) {
	value, exists := e.Attributes[key]
	return value, exists
}

// SetAttribute sets an attribute on an edge
func (e *Edge) SetAttribute(key string, value interface{}) {
	if e.Attributes == nil {
		e.Attributes = make(Attributes)
	}
	e.Attributes[key] = value
	e.UpdatedAt = time.Now()
}
