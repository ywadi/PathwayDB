package tests

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ywadi/PathwayDB/models"
	"github.com/ywadi/PathwayDB/storage"
)

// TestStorageEngine provides comprehensive tests for the storage layer
type TestStorageEngine struct {
	engine   storage.StorageEngine
	testPath string
	graphID  models.GraphID
}

// setupTestEngine creates a new test engine with a temporary database
func setupTestEngine(t *testing.T) *TestStorageEngine {
	testPath := filepath.Join(os.TempDir(), "pathwaydb_test_"+t.Name())
	engine := storage.NewBadgerEngine()
	
	err := engine.Open(testPath)
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}
	
	return &TestStorageEngine{
		engine:   engine,
		testPath: testPath,
		graphID:  models.GraphID("test-graph"),
	}
}

// cleanup closes the engine and removes test data
func (te *TestStorageEngine) cleanup() {
	if te.engine != nil {
		te.engine.Close()
	}
	os.RemoveAll(te.testPath)
}

// TestGraphCRUD tests all graph CRUD operations
func TestGraphCRUD(t *testing.T) {
	te := setupTestEngine(t)
	defer te.cleanup()

	// Test CreateGraph
	t.Run("CreateGraph", func(t *testing.T) {
		graph := &models.Graph{
			ID:          "test-graph",
			Name:        "Test Graph",
			Description: "A test graph",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		
		err := te.engine.CreateGraph(graph)
		if err != nil {
			t.Errorf("Failed to create graph: %v", err)
		}
		
		// Try to create duplicate graph (storage engine allows overwrites)
		err = te.engine.CreateGraph(graph)
		if err != nil {
			t.Errorf("Unexpected error when creating duplicate graph: %v", err)
		}
	})

	// Test GetGraph
	t.Run("GetGraph", func(t *testing.T) {
		graph, err := te.engine.GetGraph("test-graph")
		if err != nil {
			t.Errorf("Failed to get graph: %v", err)
		}
		if graph.ID != "test-graph" {
			t.Errorf("Expected graph ID 'test-graph', got %s", graph.ID)
		}
		if graph.Name != "Test Graph" {
			t.Errorf("Expected graph name 'Test Graph', got %s", graph.Name)
		}
		
		// Test non-existent graph
		_, err = te.engine.GetGraph("non-existent")
		if err == nil {
			t.Error("Expected error when getting non-existent graph")
		}
	})

	// Test UpdateGraph
	t.Run("UpdateGraph", func(t *testing.T) {
		graph, _ := te.engine.GetGraph("test-graph")
		graph.Name = "Updated Test Graph"
		graph.Description = "Updated description"
		graph.UpdatedAt = time.Now()
		
		err := te.engine.UpdateGraph(graph)
		if err != nil {
			t.Errorf("Failed to update graph: %v", err)
		}
		
		// Verify update
		updatedGraph, err := te.engine.GetGraph("test-graph")
		if err != nil {
			t.Errorf("Failed to get updated graph: %v", err)
		}
		if updatedGraph.Name != "Updated Test Graph" {
			t.Errorf("Expected updated name 'Updated Test Graph', got %s", updatedGraph.Name)
		}
	})

	// Test ListGraphs
	t.Run("ListGraphs", func(t *testing.T) {
		// Create additional graph
		graph2 := &models.Graph{
			ID:          "test-graph-2",
			Name:        "Test Graph 2",
			Description: "Second test graph",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		te.engine.CreateGraph(graph2)
		
		graphs, err := te.engine.ListGraphs()
		if err != nil {
			t.Errorf("Failed to list graphs: %v", err)
		}
		if len(graphs) < 2 {
			t.Errorf("Expected at least 2 graphs, got %d", len(graphs))
		}
	})

	// Test DeleteGraph
	t.Run("DeleteGraph", func(t *testing.T) {
		err := te.engine.DeleteGraph("test-graph-2")
		if err != nil {
			t.Errorf("Failed to delete graph: %v", err)
		}
		
		// Verify deletion
		_, err = te.engine.GetGraph("test-graph-2")
		if err == nil {
			t.Error("Expected error when getting deleted graph")
		}
		
		// Try to delete non-existent graph (storage engine handles gracefully)
		err = te.engine.DeleteGraph("non-existent")
		if err != nil {
			t.Errorf("Unexpected error when deleting non-existent graph: %v", err)
		}
	})

	t.Run("DeleteGraphWithContent", func(t *testing.T) {
		// Create a graph with content
		graphWithContentID := models.GraphID("graph-with-content")
		graphWithContent := &models.Graph{ID: graphWithContentID, Name: "Graph With Content", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		te.engine.CreateGraph(graphWithContent)

		node := &models.Node{ID: "node-in-deleted-graph", Type: "service", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		te.engine.CreateNode(graphWithContentID, node)

		// Delete the graph
		err := te.engine.DeleteGraph(graphWithContentID)
		if err != nil {
			t.Fatalf("Failed to delete graph with content: %v", err)
		}

		// Verify graph is deleted
		_, err = te.engine.GetGraph(graphWithContentID)
		if err == nil {
			t.Error("Graph 'graph-with-content' should have been deleted")
		}

		// Verify node is also deleted
		_, err = te.engine.GetNode(graphWithContentID, "node-in-deleted-graph")
		if err == nil {
			t.Error("Node in deleted graph should also be deleted")
		}
	})
}

// TestNodeCRUD tests all node CRUD operations
func TestNodeCRUD(t *testing.T) {
	te := setupTestEngine(t)
	defer te.cleanup()

	// Create test graph
	graph := &models.Graph{
		ID:          te.graphID,
		Name:        "Test Graph",
		Description: "A test graph",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	te.engine.CreateGraph(graph)

	// Test CreateNode
	t.Run("CreateNode", func(t *testing.T) {
		node := &models.Node{
			ID:   "test-node",
			Type: "service",
			Attributes: models.Attributes{
				"name":     "Test Service",
				"language": "go",
				"port":     8080,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		err := te.engine.CreateNode(te.graphID, node)
		if err != nil {
			t.Errorf("Failed to create node: %v", err)
		}
		
		// Try to create duplicate node (storage engine allows overwrites)
		err = te.engine.CreateNode(te.graphID, node)
		if err != nil {
			t.Errorf("Unexpected error when creating duplicate node: %v", err)
		}
	})

	// Test GetNode
	t.Run("GetNode", func(t *testing.T) {
		node, err := te.engine.GetNode(te.graphID, "test-node")
		if err != nil {
			t.Errorf("Failed to get node: %v", err)
		}
		if node.ID != "test-node" {
			t.Errorf("Expected node ID 'test-node', got %s", node.ID)
		}
		if node.Type != "service" {
			t.Errorf("Expected node type 'service', got %s", node.Type)
		}
		
		// Test non-existent node
		_, err = te.engine.GetNode(te.graphID, "non-existent")
		if err == nil {
			t.Error("Expected error when getting non-existent node")
		}
	})

	// Test UpdateNode
	t.Run("UpdateNode", func(t *testing.T) {
		node, _ := te.engine.GetNode(te.graphID, "test-node")
		node.SetAttribute("status", "running")
		node.SetAttribute("replicas", 3)
		node.UpdatedAt = time.Now()
		
		err := te.engine.UpdateNode(te.graphID, node)
		if err != nil {
			t.Errorf("Failed to update node: %v", err)
		}
		
		// Verify update
		updatedNode, err := te.engine.GetNode(te.graphID, "test-node")
		if err != nil {
			t.Errorf("Failed to get updated node: %v", err)
		}
		status, exists := updatedNode.GetAttribute("status")
		if !exists || status != "running" {
			t.Errorf("Expected status 'running', got %v", status)
		}
	})

	// Test ListNodes
	t.Run("ListNodes", func(t *testing.T) {
		// Create additional nodes
		nodes := []*models.Node{
			{ID: "node-1", Type: "database", Attributes: models.Attributes{"name": "DB 1"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "node-2", Type: "service", Attributes: models.Attributes{"name": "Service 2"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "node-3", Type: "cache", Attributes: models.Attributes{"name": "Cache 1"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		
		for _, node := range nodes {
			te.engine.CreateNode(te.graphID, node)
		}
		
		allNodes, err := te.engine.ListNodes(te.graphID)
		if err != nil {
			t.Errorf("Failed to list nodes: %v", err)
		}
		if len(allNodes) < 4 { // test-node + 3 new nodes
			t.Errorf("Expected at least 4 nodes, got %d", len(allNodes))
		}
	})

	// Test ListNodesByType
	t.Run("ListNodesByType", func(t *testing.T) {
		serviceNodes, err := te.engine.ListNodesByType(te.graphID, "service")
		if err != nil {
			t.Errorf("Failed to list service nodes: %v", err)
		}
		if len(serviceNodes) != 2 { // test-node + node-2
			t.Errorf("Expected 2 service nodes, got %d", len(serviceNodes))
		}
		
		databaseNodes, err := te.engine.ListNodesByType(te.graphID, "database")
		if err != nil {
			t.Errorf("Failed to list database nodes: %v", err)
		}
		if len(databaseNodes) != 1 {
			t.Errorf("Expected 1 database node, got %d", len(databaseNodes))
		}
	})

	// Test FindNodesByAttribute
	t.Run("FindNodesByAttribute", func(t *testing.T) {
		goNodes, err := te.engine.FindNodesByAttribute(te.graphID, "language", "go")
		if err != nil {
			t.Errorf("Failed to find nodes by attribute: %v", err)
		}
		if len(goNodes) != 1 {
			t.Errorf("Expected 1 Go node, got %d", len(goNodes))
		}
		
		// Test non-existent attribute
		nonExistent, err := te.engine.FindNodesByAttribute(te.graphID, "non-existent", "value")
		if err != nil {
			t.Errorf("Failed to search for non-existent attribute: %v", err)
		}
		if len(nonExistent) != 0 {
			t.Errorf("Expected 0 nodes with non-existent attribute, got %d", len(nonExistent))
		}
	})

	// Test DeleteNode
	t.Run("DeleteNode", func(t *testing.T) {
		err := te.engine.DeleteNode(te.graphID, "node-3")
		if err != nil {
			t.Errorf("Failed to delete node: %v", err)
		}
		
		// Verify deletion
		_, err = te.engine.GetNode(te.graphID, "node-3")
		if err == nil {
			t.Error("Expected error when getting deleted node")
		}
		
		// Test deleting non-existent node
		err = te.engine.DeleteNode(te.graphID, "non-existent")
		if err == nil {
			t.Error("Expected error when deleting non-existent node")
		}
	})

	t.Run("DeleteNodeWithEdges", func(t *testing.T) {
		// Create nodes and edges for cascading delete test
		nodeA := &models.Node{ID: "cascade-a", Type: "service", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		nodeB := &models.Node{ID: "cascade-b", Type: "service", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		nodeC := &models.Node{ID: "cascade-c", Type: "database", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		te.engine.CreateNode(te.graphID, nodeA)
		te.engine.CreateNode(te.graphID, nodeB)
		te.engine.CreateNode(te.graphID, nodeC)

		edgeAB := &models.Edge{ID: "edge-ab", FromNodeID: "cascade-a", ToNodeID: "cascade-b", Type: "calls", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		edgeBC := &models.Edge{ID: "edge-bc", FromNodeID: "cascade-b", ToNodeID: "cascade-c", Type: "calls", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		te.engine.CreateEdge(te.graphID, edgeAB)
		te.engine.CreateEdge(te.graphID, edgeBC)

		// Delete node B
		err := te.engine.DeleteNode(te.graphID, "cascade-b")
		if err != nil {
			t.Fatalf("Failed to delete node for cascading test: %v", err)
		}

		// Verify node B is deleted
		_, err = te.engine.GetNode(te.graphID, "cascade-b")
		if err == nil {
			t.Error("Node 'cascade-b' should have been deleted")
		}

		// Verify edges are deleted
		_, err = te.engine.GetEdge(te.graphID, "edge-ab")
		if err == nil {
			t.Error("Edge 'edge-ab' should have been deleted")
		}
		_, err = te.engine.GetEdge(te.graphID, "edge-bc")
		if err == nil {
			t.Error("Edge 'edge-bc' should have been deleted")
		}
	})
}

// TestEdgeCRUD tests all edge CRUD operations
func TestEdgeCRUD(t *testing.T) {
	te := setupTestEngine(t)
	defer te.cleanup()

	// Create test graph and nodes
	graph := &models.Graph{
		ID:          te.graphID,
		Name:        "Test Graph",
		Description: "A test graph",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	te.engine.CreateGraph(graph)

	nodes := []*models.Node{
		{ID: "node-a", Type: "service", Attributes: models.Attributes{"name": "Service A"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "node-b", Type: "service", Attributes: models.Attributes{"name": "Service B"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "node-c", Type: "database", Attributes: models.Attributes{"name": "Database C"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	
	for _, node := range nodes {
		te.engine.CreateNode(te.graphID, node)
	}

	// Test CreateEdge
	t.Run("CreateEdge", func(t *testing.T) {
		edge := &models.Edge{
			ID:         "test-edge",
			Type:       "depends_on",
			FromNodeID: "node-a",
			ToNodeID:   "node-c",
			Attributes: models.Attributes{
				"connection_type": "read_write",
				"critical":        true,
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		err := te.engine.CreateEdge(te.graphID, edge)
		if err != nil {
			t.Errorf("Failed to create edge: %v", err)
		}
		
		// Try to create duplicate edge (storage engine allows overwrites)
		err = te.engine.CreateEdge(te.graphID, edge)
		if err != nil {
			t.Errorf("Unexpected error when creating duplicate edge: %v", err)
		}
	})

	// Test GetEdge
	t.Run("GetEdge", func(t *testing.T) {
		edge, err := te.engine.GetEdge(te.graphID, "test-edge")
		if err != nil {
			t.Errorf("Failed to get edge: %v", err)
		}
		if edge.ID != "test-edge" {
			t.Errorf("Expected edge ID 'test-edge', got %s", edge.ID)
		}
		if edge.FromNodeID != "node-a" {
			t.Errorf("Expected from node 'node-a', got %s", edge.FromNodeID)
		}
		if edge.ToNodeID != "node-c" {
			t.Errorf("Expected to node 'node-c', got %s", edge.ToNodeID)
		}
		
		// Test non-existent edge
		_, err = te.engine.GetEdge(te.graphID, "non-existent")
		if err == nil {
			t.Error("Expected error when getting non-existent edge")
		}
	})

	// Test UpdateEdge
	t.Run("UpdateEdge", func(t *testing.T) {
		edge, _ := te.engine.GetEdge(te.graphID, "test-edge")
		edge.SetAttribute("priority", "high")
		edge.SetAttribute("timeout", 30)
		edge.UpdatedAt = time.Now()
		
		err := te.engine.UpdateEdge(te.graphID, edge)
		if err != nil {
			t.Errorf("Failed to update edge: %v", err)
		}
		
		// Verify update
		updatedEdge, err := te.engine.GetEdge(te.graphID, "test-edge")
		if err != nil {
			t.Errorf("Failed to get updated edge: %v", err)
		}
		priority, exists := updatedEdge.GetAttribute("priority")
		if !exists || priority != "high" {
			t.Errorf("Expected priority 'high', got %v", priority)
		}
	})

	// Test ListEdges
	t.Run("ListEdges", func(t *testing.T) {
		// Create additional edges
		edges := []*models.Edge{
			{ID: "edge-1", Type: "depends_on", FromNodeID: "node-b", ToNodeID: "node-c", Attributes: models.Attributes{}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: "edge-2", Type: "calls", FromNodeID: "node-a", ToNodeID: "node-b", Attributes: models.Attributes{}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}
		
		for _, edge := range edges {
			te.engine.CreateEdge(te.graphID, edge)
		}
		
		allEdges, err := te.engine.ListEdges(te.graphID)
		if err != nil {
			t.Errorf("Failed to list edges: %v", err)
		}
		if len(allEdges) != 3 { // test-edge + 2 new edges
			t.Errorf("Expected 3 edges, got %d", len(allEdges))
		}
	})

	// Test ListEdgesByType
	t.Run("ListEdgesByType", func(t *testing.T) {
		dependsOnEdges, err := te.engine.ListEdgesByType(te.graphID, "depends_on")
		if err != nil {
			t.Errorf("Failed to list depends_on edges: %v", err)
		}
		if len(dependsOnEdges) != 2 {
			t.Errorf("Expected 2 depends_on edges, got %d", len(dependsOnEdges))
		}
		
		callsEdges, err := te.engine.ListEdgesByType(te.graphID, "calls")
		if err != nil {
			t.Errorf("Failed to list calls edges: %v", err)
		}
		if len(callsEdges) != 1 {
			t.Errorf("Expected 1 calls edge, got %d", len(callsEdges))
		}
	})

	// Test relationship queries
	t.Run("RelationshipQueries", func(t *testing.T) {
		// Test GetOutgoingEdges
		outgoingFromA, err := te.engine.GetOutgoingEdges(te.graphID, "node-a")
		if err != nil {
			t.Errorf("Failed to get outgoing edges from node-a: %v", err)
		}
		if len(outgoingFromA) != 2 { // test-edge + edge-2
			t.Errorf("Expected 2 outgoing edges from node-a, got %d", len(outgoingFromA))
		}
		
		// Test GetIncomingEdges
		incomingToC, err := te.engine.GetIncomingEdges(te.graphID, "node-c")
		if err != nil {
			t.Errorf("Failed to get incoming edges to node-c: %v", err)
		}
		if len(incomingToC) != 2 { // test-edge + edge-1
			t.Errorf("Expected 2 incoming edges to node-c, got %d", len(incomingToC))
		}
		
		// Test GetConnectedNodes
		connectedToA, err := te.engine.GetConnectedNodes(te.graphID, "node-a")
		if err != nil {
			t.Errorf("Failed to get connected nodes to node-a: %v", err)
		}
		if len(connectedToA) != 2 { // node-b and node-c
			t.Errorf("Expected 2 connected nodes to node-a, got %d", len(connectedToA))
		}

		// Test with node IDs containing colons
		nodeWithColon1 := &models.Node{ID: "user:1", Type: "user", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		nodeWithColon2 := &models.Node{ID: "product:123", Type: "product", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		te.engine.CreateNode(te.graphID, nodeWithColon1)
		te.engine.CreateNode(te.graphID, nodeWithColon2)
		edgeWithColon := &models.Edge{ID: "user:1-buys-product:123", Type: "buys", FromNodeID: "user:1", ToNodeID: "product:123", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		te.engine.CreateEdge(te.graphID, edgeWithColon)

		outgoingFromUser, err := te.engine.GetOutgoingEdges(te.graphID, "user:1")
		if err != nil {
			t.Errorf("Failed to get outgoing edges from node 'user:1': %v", err)
		}
		if len(outgoingFromUser) != 1 {
			t.Errorf("Expected 1 outgoing edge from 'user:1', got %d", len(outgoingFromUser))
		} else if outgoingFromUser[0].ID != "user:1-buys-product:123" {
			t.Errorf("Incorrect edge ID found for 'user:1', got %s", outgoingFromUser[0].ID)
		}
	})

	// Test DeleteEdge
	t.Run("DeleteEdge", func(t *testing.T) {
		err := te.engine.DeleteEdge(te.graphID, "edge-2")
		if err != nil {
			t.Errorf("Failed to delete edge: %v", err)
		}
		
		// Verify deletion
		_, err = te.engine.GetEdge(te.graphID, "edge-2")
		if err == nil {
			t.Error("Expected error when getting deleted edge")
		}
		
		// Test deleting non-existent edge
		err = te.engine.DeleteEdge(te.graphID, "non-existent")
		if err == nil {
			t.Error("Expected error when deleting non-existent edge")
		}
	})
}

// TestDatabaseOperations tests database-level operations
func TestDatabaseOperations(t *testing.T) {
	te := setupTestEngine(t)
	defer te.cleanup()

	// Test database is opened
	t.Run("DatabaseOpen", func(t *testing.T) {
		// Try to create a graph to verify database is working
		graph := &models.Graph{
			ID:          te.graphID,
			Name:        "Test Graph",
			Description: "A test graph",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		
		err := te.engine.CreateGraph(graph)
		if err != nil {
			t.Errorf("Database not properly opened: %v", err)
		}
	})

	// Test backup functionality
	t.Run("Backup", func(t *testing.T) {
		// Create backup directory first
		backupDir := filepath.Join("/tmp", "pathwaydb_backup_test")
		err := os.MkdirAll(backupDir, 0755)
		if err != nil {
			t.Errorf("Failed to create backup directory: %v", err)
		}
		defer os.RemoveAll(backupDir)
		
		// Create backup (pass directory, engine will create backup.db inside)
		err = te.engine.Backup(backupDir)
		if err != nil {
			t.Errorf("Failed to create backup: %v", err)
		}
		
		// Verify backup file exists
		backupPath := filepath.Join(backupDir, "backup.db")
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			t.Error("Backup file was not created")
		}
	})

	// Test transaction functionality is handled internally by the storage engine
	t.Run("TransactionSupport", func(t *testing.T) {
		// Test that operations are atomic (create multiple items)
		graph := &models.Graph{
			ID:          "tx-test",
			Name:        "Transaction Test",
			Description: "Test atomic operations",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		
		err := te.engine.CreateGraph(graph)
		if err != nil {
			t.Errorf("Failed to create graph for transaction test: %v", err)
		}
		
		// Create node and verify it exists (tests internal transaction handling)
		node := &models.Node{
			ID:         "tx-node",
			Type:       "service",
			Attributes: models.Attributes{"name": "TX Node"},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		
		err = te.engine.CreateNode(graph.ID, node)
		if err != nil {
			t.Errorf("Failed to create node in transaction: %v", err)
		}
		
		// Verify node exists
		retrievedNode, err := te.engine.GetNode(graph.ID, node.ID)
		if err != nil {
			t.Errorf("Node not found after transaction: %v", err)
		}
		if retrievedNode.ID != node.ID {
			t.Errorf("Retrieved node ID mismatch: expected %s, got %s", node.ID, retrievedNode.ID)
		}
	})
}

// TestErrorHandling tests various error conditions
func TestTTL(t *testing.T) {
	te := setupTestEngine(t)
	defer te.cleanup()

	graph := &models.Graph{ID: te.graphID, Name: "TTL Test Graph"}
	te.engine.CreateGraph(graph)

	t.Run("NodeExpiration", func(t *testing.T) {
		nodeID := models.NodeID("ttl-node-1")
		expiresAt := time.Now().Add(2 * time.Second)
		node := &models.Node{ID: nodeID, Type: "service", ExpiresAt: &expiresAt}

		if err := te.engine.CreateNode(te.graphID, node); err != nil {
			t.Fatalf("Failed to create node with TTL: %v", err)
		}

		// Check it exists
		_, err := te.engine.GetNode(te.graphID, nodeID)
		if err != nil {
			t.Fatalf("Node should exist immediately after creation: %v", err)
		}

		// Wait for expiration
		time.Sleep(3 * time.Second)

		// Manually trigger cleanup for test reliability
		type ttlManager interface {
			Cleanup()
		}
		if manager, ok := te.engine.(ttlManager); ok {
			manager.Cleanup()
		}
		time.Sleep(1 * time.Second) // Give cleanup a moment

		// Check it's gone
		_, err = te.engine.GetNode(te.graphID, nodeID)
		if err == nil {
			t.Fatal("Node should be deleted after TTL expires")
		}
	})

	t.Run("EdgeExpiration", func(t *testing.T) {
		nodeC := &models.Node{ID: "node-c", Type: "service"}
		nodeD := &models.Node{ID: "node-d", Type: "service"}
		te.engine.CreateNode(te.graphID, nodeC)
		te.engine.CreateNode(te.graphID, nodeD)

		edgeID := models.EdgeID("ttl-edge-1")
		expiresAt := time.Now().Add(2 * time.Second)
		edge := &models.Edge{ID: edgeID, FromNodeID: "node-c", ToNodeID: "node-d", Type: "calls", ExpiresAt: &expiresAt}

		if err := te.engine.CreateEdge(te.graphID, edge); err != nil {
			t.Fatalf("Failed to create edge with TTL: %v", err)
		}

		// Check it exists
		_, err := te.engine.GetEdge(te.graphID, edgeID)
		if err != nil {
			t.Fatalf("Edge should exist immediately after creation: %v", err)
		}

		// Wait for expiration
		time.Sleep(3 * time.Second)

		// Check it's gone
		_, err = te.engine.GetEdge(te.graphID, edgeID)
		if err == nil {
			t.Fatal("Edge should be deleted after TTL expires")
		}

		// Check that nodes still exist
		_, err = te.engine.GetNode(te.graphID, "node-c")
		if err != nil {
			t.Fatal("Node C should still exist")
		}
		_, err = te.engine.GetNode(te.graphID, "node-d")
		if err != nil {
			t.Fatal("Node D should still exist")
		}
	})

	t.Run("NodeCascadingDelete", func(t *testing.T) {
		nodeE := &models.Node{ID: "node-e", Type: "service"}
		nodeF := &models.Node{ID: "node-f", Type: "service"}
		expiresAt := time.Now().Add(2 * time.Second)
		nodeE.ExpiresAt = &expiresAt

		te.engine.CreateNode(te.graphID, nodeE)
		te.engine.CreateNode(te.graphID, nodeF)

		edgeEF := &models.Edge{ID: "edge-ef", FromNodeID: "node-e", ToNodeID: "node-f", Type: "calls"}
		te.engine.CreateEdge(te.graphID, edgeEF)

		// Wait for expiration
		time.Sleep(3 * time.Second)

		// Manually trigger cleanup for test reliability
		type ttlManager interface {
			Cleanup()
		}
		if manager, ok := te.engine.(ttlManager); ok {
			manager.Cleanup()
		}
		time.Sleep(1 * time.Second) // Give cleanup a moment

		// Check node E is gone
		_, err := te.engine.GetNode(te.graphID, "node-e")
		if err == nil {
			t.Fatal("Node E should be deleted after TTL expires")
		}

		// Check edge EF is gone
		_, err = te.engine.GetEdge(te.graphID, "edge-ef")
		if err == nil {
			t.Fatal("Edge EF should be deleted due to cascading delete")
		}

		// Check node F still exists
		_, err = te.engine.GetNode(te.graphID, "node-f")
		if err != nil {
			t.Fatal("Node F should still exist")
		}
	})
}

func TestErrorHandling(t *testing.T) {
	// Test operations on closed database
	t.Run("ClosedDatabase", func(t *testing.T) {
		engine := storage.NewBadgerEngine()
		// Don't open the database
		
		graph := &models.Graph{ID: "test", Name: "Test", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		err := engine.CreateGraph(graph)
		if err == nil {
			t.Error("Expected error when operating on closed database")
		}
	})

	// Test operations on non-existent graph
	t.Run("NonExistentGraph", func(t *testing.T) {
		te := setupTestEngine(t)
		defer te.cleanup()
		
		node := &models.Node{ID: "test", Type: "test", Attributes: models.Attributes{}, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		// Try to create node in non-existent graph (storage engine handles gracefully)
		err := te.engine.CreateNode("non-existent-graph", node)
		if err != nil {
			t.Errorf("Unexpected error when creating node in non-existent graph: %v", err)
		}
	})

	// Test invalid edge creation (non-existent nodes)
	t.Run("InvalidEdge", func(t *testing.T) {
		te := setupTestEngine(t)
		defer te.cleanup()
		
		// Create graph but no nodes
		graph := &models.Graph{ID: "test-graph", Name: "Test", CreatedAt: time.Now(), UpdatedAt: time.Now()}
		te.engine.CreateGraph(graph)
		
		edge := &models.Edge{
			ID:         "test-edge",
			Type:       "depends_on",
			FromNodeID: "non-existent-from",
			ToNodeID:   "non-existent-to",
			Attributes: models.Attributes{},
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		
		err := te.engine.CreateEdge(graph.ID, edge)
		if err == nil {
			t.Error("Expected error when creating edge with non-existent nodes")
		}
	})
}
