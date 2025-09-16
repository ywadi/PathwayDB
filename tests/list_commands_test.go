package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ywadi/PathwayDB/models"
	"github.com/ywadi/PathwayDB/redis/commands"
	"github.com/ywadi/PathwayDB/storage"
)

// TestListCommands tests NODE.LIST and EDGE.LIST output formats
func TestListCommands(t *testing.T) {
	testPath := filepath.Join(os.TempDir(), "pathwaydb_list_test")
	engine := storage.NewBadgerEngine()
	defer func() {
		engine.Close()
		os.RemoveAll(testPath)
	}()

	err := engine.Open(testPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	graphID := models.GraphID("list-test-graph")

	// Create test graph
	graph := &models.Graph{
		ID:          graphID,
		Name:        "List Test Graph",
		Description: "Graph for testing list commands",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = engine.CreateGraph(graph)
	if err != nil {
		t.Fatalf("Failed to create graph: %v", err)
	}

	// Create test nodes
	nodes := []*models.Node{
		{ID: "service-a", Type: "service", Attributes: models.Attributes{"name": "Service A"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "service-b", Type: "service", Attributes: models.Attributes{"name": "Service B"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "database-1", Type: "database", Attributes: models.Attributes{"name": "Database 1"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for _, node := range nodes {
		err = engine.CreateNode(graphID, node)
		if err != nil {
			t.Fatalf("Failed to create node %s: %v", node.ID, err)
		}
	}

	// Create test edges
	edges := []*models.Edge{
		{ID: "edge-ab", Type: "depends_on", FromNodeID: "service-a", ToNodeID: "service-b", Attributes: models.Attributes{"type": "http"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "edge-a-db", Type: "connects_to", FromNodeID: "service-a", ToNodeID: "database-1", Attributes: models.Attributes{"type": "sql"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for _, edge := range edges {
		err = engine.CreateEdge(graphID, edge)
		if err != nil {
			t.Fatalf("Failed to create edge %s: %v", edge.ID, err)
		}
	}

	// Test NODE.LIST command
	t.Run("NodeListFormat", func(t *testing.T) {
		nodeCommands := commands.NewNodeCommands(engine)
		response, err := nodeCommands.Handle("LIST", []string{string(graphID)})
		if err != nil {
			t.Fatalf("NODE.LIST failed: %v", err)
		}

		if response.Type != 2 { // ResponseTypeArray = 2
			t.Fatalf("Expected array response, got %d", response.Type)
		}

		items := response.ArrayValue
		if len(items) != 3 {
			t.Fatalf("Expected 3 items, got %d", len(items))
		}

		// Verify id:type format - check that all expected items are present
		expectedItems := []string{"service-a:service", "service-b:service", "database-1:database"}
		for _, expected := range expectedItems {
			found := false
			for _, item := range items {
				if item == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected to find %s in items %v", expected, items)
			}
		}

		// Verify all items contain colon separator
		for i, item := range items {
			if !strings.Contains(item, ":") {
				t.Errorf("Item %d (%s) should contain colon separator for id:type format", i, item)
			}
		}
	})

	// Test EDGE.LIST command
	t.Run("EdgeListFormat", func(t *testing.T) {
		edgeCommands := commands.NewEdgeCommands(engine)
		response, err := edgeCommands.Handle("LIST", []string{string(graphID)})
		if err != nil {
			t.Fatalf("EDGE.LIST failed: %v", err)
		}

		if response.Type != 2 { // ResponseTypeArray = 2
			t.Fatalf("Expected array response, got %d", response.Type)
		}

		items := response.ArrayValue
		if len(items) != 2 {
			t.Fatalf("Expected 2 items, got %d", len(items))
		}

		// Verify id:type format - check that all expected items are present
		expectedItems := []string{"edge-ab:depends_on", "edge-a-db:connects_to"}
		for _, expected := range expectedItems {
			found := false
			for _, item := range items {
				if item == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected to find %s in items %v", expected, items)
			}
		}

		// Verify all items contain colon separator
		for i, item := range items {
			if !strings.Contains(item, ":") {
				t.Errorf("Item %d (%s) should contain colon separator for id:type format", i, item)
			}
		}
	})
}

// TestEdgeNeighborsCommand tests EDGE.NEIGHBORS output format with arrow notation
func TestEdgeNeighborsCommand(t *testing.T) {
	testPath := filepath.Join(os.TempDir(), "pathwaydb_neighbors_test")
	engine := storage.NewBadgerEngine()
	defer func() {
		engine.Close()
		os.RemoveAll(testPath)
	}()

	err := engine.Open(testPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}

	graphID := models.GraphID("neighbors-test-graph")

	// Create test graph
	graph := &models.Graph{
		ID:          graphID,
		Name:        "Neighbors Test Graph",
		Description: "Graph for testing edge neighbors command",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = engine.CreateGraph(graph)
	if err != nil {
		t.Fatalf("Failed to create graph: %v", err)
	}

	// Create test nodes
	nodes := []*models.Node{
		{ID: "node-a", Type: "service", Attributes: models.Attributes{"name": "Node A"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "node-b", Type: "service", Attributes: models.Attributes{"name": "Node B"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "node-c", Type: "database", Attributes: models.Attributes{"name": "Node C"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for _, node := range nodes {
		err = engine.CreateNode(graphID, node)
		if err != nil {
			t.Fatalf("Failed to create node %s: %v", node.ID, err)
		}
	}

	// Create test edges
	edges := []*models.Edge{
		{ID: "edge-ab", Type: "depends_on", FromNodeID: "node-a", ToNodeID: "node-b", Attributes: models.Attributes{"type": "http"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "edge-ac", Type: "connects_to", FromNodeID: "node-a", ToNodeID: "node-c", Attributes: models.Attributes{"type": "sql"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "edge-ba", Type: "notifies", FromNodeID: "node-b", ToNodeID: "node-a", Attributes: models.Attributes{"type": "async"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for _, edge := range edges {
		err = engine.CreateEdge(graphID, edge)
		if err != nil {
			t.Fatalf("Failed to create edge %s: %v", edge.ID, err)
		}
	}

	// Test EDGE.NEIGHBORS command for outgoing edges
	t.Run("EdgeNeighborsOutgoing", func(t *testing.T) {
		edgeCommands := commands.NewEdgeCommands(engine)
		response, err := edgeCommands.Handle("NEIGHBORS", []string{string(graphID), "node-a", "out"})
		if err != nil {
			t.Fatalf("EDGE.NEIGHBORS failed: %v", err)
		}

		if response.Type != 2 { // ResponseTypeArray = 2
			t.Fatalf("Expected array response, got %d", response.Type)
		}

		items := response.ArrayValue
		if len(items) < 3 { // Count + at least 2 neighbors
			t.Fatalf("Expected at least 3 items (count + 2 neighbors), got %d", len(items))
		}

		// First item should be count
		count := items[0]
		if count != "2" {
			t.Errorf("Expected count to be 2, got %s", count)
		}

		// Verify arrow notation format for neighbor items (skip count)
		neighborItems := items[1:]
		for i, item := range neighborItems {
			if !strings.Contains(item, "->") {
				t.Errorf("Outgoing neighbor %d (%s) should contain '->' arrow", i, item)
			}
			if strings.Contains(item, "<-") {
				t.Errorf("Outgoing neighbor %s should not contain '<-' arrow", item)
			}
		}

		// Check specific format: neighbor_node:type->edge:type
		expectedPatterns := []string{
			"node-b:service->edge-ab:depends_on",
			"node-c:database->edge-ac:connects_to",
		}

		for _, expected := range expectedPatterns {
			found := false
			for _, item := range neighborItems {
				if item == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected to find pattern %s in neighbors output %v", expected, neighborItems)
			}
		}
	})

	// Test EDGE.NEIGHBORS command for incoming edges
	t.Run("EdgeNeighborsIncoming", func(t *testing.T) {
		edgeCommands := commands.NewEdgeCommands(engine)
		response, err := edgeCommands.Handle("NEIGHBORS", []string{string(graphID), "node-a", "in"})
		if err != nil {
			t.Fatalf("EDGE.NEIGHBORS failed: %v", err)
		}

		if response.Type != 2 { // ResponseTypeArray = 2
			t.Fatalf("Expected array response, got %d", response.Type)
		}

		items := response.ArrayValue
		if len(items) < 2 { // Count + at least 1 neighbor
			t.Fatalf("Expected at least 2 items (count + 1 neighbor), got %d", len(items))
		}

		// First item should be count
		count := items[0]
		if count != "1" {
			t.Errorf("Expected count to be 1, got %s", count)
		}

		// Verify arrow notation format for incoming edges (skip count)
		neighborItems := items[1:]
		incomingFound := false
		for _, item := range neighborItems {
			if strings.Contains(item, "<-") {
				incomingFound = true
			}
		}
		if !incomingFound {
			t.Errorf("Expected to find at least one incoming neighbor with '<-' arrow in %v", neighborItems)
		}
	})

	// Test EDGE.NEIGHBORS command for both directions
	t.Run("EdgeNeighborsBoth", func(t *testing.T) {
		edgeCommands := commands.NewEdgeCommands(engine)
		response, err := edgeCommands.Handle("NEIGHBORS", []string{string(graphID), "node-a", "both"})
		if err != nil {
			t.Fatalf("EDGE.NEIGHBORS failed: %v", err)
		}

		if response.Type != 2 { // ResponseTypeArray = 2
			t.Fatalf("Expected array response, got %d", response.Type)
		}

		items := response.ArrayValue
		if len(items) < 4 { // Count + at least 3 neighbors
			t.Fatalf("Expected at least 4 items (count + 3 neighbors), got %d", len(items))
		}

		// First item should be count
		count := items[0]
		if count != "3" {
			t.Errorf("Expected count to be 3, got %s", count)
		}

		// Count arrows to verify mix of directions (skip count)
		neighborItems := items[1:]
		outgoingCount := 0
		incomingCount := 0
		for _, item := range neighborItems {
			if strings.Contains(item, "->") && !strings.Contains(item, "<-") {
				outgoingCount++
			} else if strings.Contains(item, "<-") {
				incomingCount++
			}
		}

		if outgoingCount != 2 {
			t.Errorf("Expected 2 outgoing neighbors, got %d", outgoingCount)
		}
		if incomingCount != 1 {
			t.Errorf("Expected 1 incoming neighbor, got %d", incomingCount)
		}
	})
}
