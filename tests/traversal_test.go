package tests

import (
	"testing"

	"github.com/ywadi/PathwayDB/models"
	"github.com/ywadi/PathwayDB/types"
)

// TestTraversal tests the full graph traversal functionality
func TestTraversal(t *testing.T) {
	te := setupTestAnalysisEngine(t)
	defer te.cleanup()

	// Setup a graph for traversal tests
	nodes := []*models.Node{
		{ID: "a", Type: "service"},
		{ID: "b", Type: "service"},
		{ID: "c", Type: "database"},
		{ID: "d", Type: "cache"},
		{ID: "e", Type: "service"},
	}
	for _, node := range nodes {
		if err := te.engine.CreateNode(te.graphID, node); err != nil {
			t.Fatalf("Failed to create node: %v", err)
		}
	}

	edges := []*models.Edge{
		{ID: "ab", FromNodeID: "a", ToNodeID: "b", Type: "calls"},
		{ID: "bc", FromNodeID: "b", ToNodeID: "c", Type: "writes_to"},
		{ID: "bd", FromNodeID: "b", ToNodeID: "d", Type: "uses"},
		{ID: "ae", FromNodeID: "a", ToNodeID: "e", Type: "calls"},
	}
	for _, edge := range edges {
		if err := te.engine.CreateEdge(te.graphID, edge); err != nil {
			t.Fatalf("Failed to create edge: %v", err)
		}
	}

	t.Run("ForwardTraversal", func(t *testing.T) {
		result, err := te.analyzer.DepthFirstSearch(te.graphID, "a", &types.TraversalOptions{Direction: types.DirectionForward, MaxDepth: -1})
		if err != nil {
			t.Fatalf("Forward traversal failed: %v", err)
		}
		if len(result.Nodes) != 5 { // a, b, c, d, e
			t.Errorf("Expected 5 nodes in forward traversal, got %d", len(result.Nodes))
		}
	})

	t.Run("BackwardTraversal", func(t *testing.T) {
		result, err := te.analyzer.DepthFirstSearch(te.graphID, "c", &types.TraversalOptions{Direction: types.DirectionBackward, MaxDepth: -1})
		if err != nil {
			t.Fatalf("Backward traversal failed: %v", err)
		}
		if len(result.Nodes) != 3 { // c, b, a
			t.Errorf("Expected 3 nodes in backward traversal from c, got %d", len(result.Nodes))
		}
	})

	t.Run("FilterByEdgeType", func(t *testing.T) {
		options := &types.TraversalOptions{Direction: types.DirectionForward, MaxDepth: -1, EdgeTypes: []models.EdgeType{"calls"}}
		result, err := te.analyzer.DepthFirstSearch(te.graphID, "a", options)
		if err != nil {
			t.Fatalf("Traversal with edge filter failed: %v", err)
		}
		if len(result.Nodes) != 3 { // a, b, e
			t.Errorf("Expected 3 nodes when filtering by 'calls' edge type, got %d", len(result.Nodes))
		}
	})

	t.Run("FilterByNodeType", func(t *testing.T) {
		options := &types.TraversalOptions{Direction: types.DirectionForward, MaxDepth: -1, NodeTypes: []models.NodeType{"service"}}
		result, err := te.analyzer.DepthFirstSearch(te.graphID, "a", options)
		if err != nil {
			t.Fatalf("Traversal with node filter failed: %v", err)
		}
		// Traversal still visits all nodes, but only returns those matching the type
		if len(result.Nodes) != 3 { // a, b, e
			t.Errorf("Expected 3 nodes when filtering by 'service' node type, got %d", len(result.Nodes))
		}
	})

	t.Run("AllPathsTraversal", func(t *testing.T) {
		allPaths, err := te.analyzer.AllPathsTraversal(te.graphID, "a", &types.TraversalOptions{Direction: types.DirectionForward, MaxDepth: -1})
		if err != nil {
			t.Fatalf("All paths traversal failed: %v", err)
		}
		// Should find multiple paths from 'a': a->b->c, a->b->d, a->e
		if len(allPaths) < 2 {
			t.Errorf("Expected at least 2 paths from node 'a', got %d", len(allPaths))
		}
		
		// Verify each path is valid
		for i, path := range allPaths {
			if len(path.Nodes) == 0 {
				t.Errorf("Path %d is empty", i)
			}
			if path.Nodes[0].ID != "a" {
				t.Errorf("Path %d should start with node 'a', got %s", i, path.Nodes[0].ID)
			}
		}
	})
}
