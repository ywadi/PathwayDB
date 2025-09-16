package storage

import (
	"github.com/ywadi/PathwayDB/models"
)

// StorageEngine defines the interface for the storage layer
type StorageEngine interface {
	// Graph operations
	CreateGraph(graph *models.Graph) error
	GetGraph(graphID models.GraphID) (*models.Graph, error)
	UpdateGraph(graph *models.Graph) error
	DeleteGraph(graphID models.GraphID) error
	ListGraphs() ([]*models.Graph, error)
	CountNodes(graphID models.GraphID) (int, error)
	CountEdges(graphID models.GraphID) (int, error)

	// Node operations
	CreateNode(graphID models.GraphID, node *models.Node) error
	GetNode(graphID models.GraphID, nodeID models.NodeID) (*models.Node, error)
	UpdateNode(graphID models.GraphID, node *models.Node) error
	DeleteNode(graphID models.GraphID, nodeID models.NodeID) error
	ListNodes(graphID models.GraphID) ([]*models.Node, error)
	ListNodesByType(graphID models.GraphID, nodeType models.NodeType) ([]*models.Node, error)

	// Edge operations
	CreateEdge(graphID models.GraphID, edge *models.Edge) error
	GetEdge(graphID models.GraphID, edgeID models.EdgeID) (*models.Edge, error)
	UpdateEdge(graphID models.GraphID, edge *models.Edge) error
	DeleteEdge(graphID models.GraphID, edgeID models.EdgeID) error
	ListEdges(graphID models.GraphID) ([]*models.Edge, error)
	ListEdgesByType(graphID models.GraphID, edgeType models.EdgeType) ([]*models.Edge, error)

	// Relationship operations
	GetOutgoingEdges(graphID models.GraphID, nodeID models.NodeID) ([]*models.Edge, error)
	GetIncomingEdges(graphID models.GraphID, nodeID models.NodeID) ([]*models.Edge, error)
	GetConnectedNodes(graphID models.GraphID, nodeID models.NodeID) ([]*models.Node, error)

	// Attribute filtering
	FindNodesByAttribute(graphID models.GraphID, attrKey string, attrValue interface{}) ([]*models.Node, error)
	FindEdgesByAttribute(graphID models.GraphID, attrKey string, attrValue interface{}) ([]*models.Edge, error)

	// Database lifecycle
	Open(path string) error
	Close() error
	Backup(path string) error
}

// FilterOptions represents options for filtering nodes/edges
type FilterOptions struct {
	NodeTypes []models.NodeType
	EdgeTypes []models.EdgeType
	Attributes map[string]interface{}
	Limit     int
	Offset    int
}

// TransactionFunc represents a function that can be executed within a transaction
type TransactionFunc func(txn Transaction) error

// Transaction defines the interface for database transactions
type Transaction interface {
	// Node operations within transaction
	CreateNode(graphID models.GraphID, node *models.Node) error
	GetNode(graphID models.GraphID, nodeID models.NodeID) (*models.Node, error)
	UpdateNode(graphID models.GraphID, node *models.Node) error
	DeleteNode(graphID models.GraphID, nodeID models.NodeID) error

	// Edge operations within transaction
	CreateEdge(graphID models.GraphID, edge *models.Edge) error
	GetEdge(graphID models.GraphID, edgeID models.EdgeID) (*models.Edge, error)
	UpdateEdge(graphID models.GraphID, edge *models.Edge) error
	DeleteEdge(graphID models.GraphID, edgeID models.EdgeID) error

	// Transaction control
	Commit() error
	Discard()
}
