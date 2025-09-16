package types

import (
	"github.com/ywadi/PathwayDB/models"
)

// TraversalResult represents the result of a graph traversal
type TraversalResult struct {
	Nodes    []*models.Node `json:"nodes"`
	Edges    []*models.Edge `json:"edges"`
	Path     []models.NodeID `json:"path"`
	Distance int            `json:"distance"`
}

// CycleResult represents a detected cycle in the graph
type CycleResult struct {
	Nodes []models.NodeID `json:"nodes"`
	Edges []models.EdgeID `json:"edges"`
}

// GraphStats represents statistics about a graph
type GraphStats struct {
	NodeCount          int                        `json:"node_count"`
	EdgeCount          int                        `json:"edge_count"`
	NodeTypeCount      map[models.NodeType]int    `json:"node_type_count"`
	EdgeTypeCount      map[models.EdgeType]int    `json:"edge_type_count"`
	MaxDepth           int                        `json:"max_depth"`
	RootNodeCount      int                        `json:"root_node_count"`
	LeafNodeCount      int                        `json:"leaf_node_count"`
	OrphanNodeCount    int                        `json:"orphan_node_count"`
	HasCycles          bool                       `json:"has_cycles"`
	ConnectedComponents int                       `json:"connected_components"`
}

// NodeMetrics represents metrics for a specific node
type NodeMetrics struct {
	NodeID           models.NodeID `json:"node_id"`
	InDegree         int           `json:"in_degree"`
	OutDegree        int           `json:"out_degree"`
	Depth            int           `json:"depth"`
	DependencyCount  int           `json:"dependency_count"`
	DependentCount   int           `json:"dependent_count"`
	IsRoot           bool          `json:"is_root"`
	IsLeaf           bool          `json:"is_leaf"`
	IsOrphan         bool          `json:"is_orphan"`
	IsCriticalPath   bool          `json:"is_critical_path"`
}

// PathResult represents a path between two nodes
type PathResult struct {
	FromNodeID models.NodeID   `json:"from_node_id"`
	ToNodeID   models.NodeID   `json:"to_node_id"`
	Path       []models.NodeID `json:"path"`
	Length     int             `json:"length"`
	Edges      []models.EdgeID `json:"edges"`
}

// DependencyTree represents a hierarchical dependency structure
type DependencyTree struct {
	NodeID   models.NodeID    `json:"node_id"`
	Node     *models.Node     `json:"node"`
	Children []*DependencyTree `json:"children"`
	Depth    int              `json:"depth"`
}

// TraversalOptions provides options for graph traversal
type TraversalOptions struct {
	MaxDepth     int                        `json:"max_depth"`
	EdgeTypes    []models.EdgeType          `json:"edge_types"`
	NodeTypes    []models.NodeType          `json:"node_types"`
	Direction    TraversalDirection         `json:"direction"`
	StopCondition func(*models.Node) bool    `json:"-"`
}

// TraversalDirection specifies the direction of traversal
type TraversalDirection int

const (
	DirectionForward TraversalDirection = iota  // Follow outgoing edges
	DirectionBackward                           // Follow incoming edges  
	DirectionBoth                              // Follow both directions
)
