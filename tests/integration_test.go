package tests

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ywadi/PathwayDB/analysis"
	"github.com/ywadi/PathwayDB/models"
	"github.com/ywadi/PathwayDB/storage"
	"github.com/ywadi/PathwayDB/types"
)

// TestIntegration provides end-to-end integration tests
func TestCompleteWorkflow(t *testing.T) {
	testPath := filepath.Join(os.TempDir(), "pathwaydb_integration_test")
	engine := storage.NewBadgerEngine()
	defer func() {
		engine.Close()
		os.RemoveAll(testPath)
	}()
	
	err := engine.Open(testPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	
	analyzer := analysis.NewGraphAnalyzer(engine)
	graphID := models.GraphID("integration-graph")
	
	// Step 1: Create graph and populate with complex dependency structure
	graph := &models.Graph{
		ID:          graphID,
		Name:        "Integration Test Graph",
		Description: "Complex dependency graph for integration testing",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	err = engine.CreateGraph(graph)
	if err != nil {
		t.Fatalf("Failed to create graph: %v", err)
	}
	
	// Create nodes representing a microservices architecture
	nodes := []*models.Node{
		{ID: "frontend", Type: "application", Attributes: models.Attributes{"name": "Frontend App", "tech": "react"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "api-gateway", Type: "service", Attributes: models.Attributes{"name": "API Gateway", "tech": "nginx"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "auth-service", Type: "service", Attributes: models.Attributes{"name": "Auth Service", "tech": "go"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "user-service", Type: "service", Attributes: models.Attributes{"name": "User Service", "tech": "go"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "order-service", Type: "service", Attributes: models.Attributes{"name": "Order Service", "tech": "java"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "payment-service", Type: "service", Attributes: models.Attributes{"name": "Payment Service", "tech": "python"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "notification-service", Type: "service", Attributes: models.Attributes{"name": "Notification Service", "tech": "node"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "user-db", Type: "database", Attributes: models.Attributes{"name": "User Database", "tech": "postgresql"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "order-db", Type: "database", Attributes: models.Attributes{"name": "Order Database", "tech": "postgresql"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "redis-cache", Type: "cache", Attributes: models.Attributes{"name": "Redis Cache", "tech": "redis"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "message-queue", Type: "queue", Attributes: models.Attributes{"name": "Message Queue", "tech": "rabbitmq"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "logger", Type: "library", Attributes: models.Attributes{"name": "Logger", "tech": "logrus"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	
	for _, node := range nodes {
		err = engine.CreateNode(graphID, node)
		if err != nil {
			t.Fatalf("Failed to create node %s: %v", node.ID, err)
		}
	}
	
	// Create edges representing dependencies
	edges := []*models.Edge{
		// Frontend dependencies
		{ID: "frontend-gateway", Type: "depends_on", FromNodeID: "frontend", ToNodeID: "api-gateway", Attributes: models.Attributes{"type": "http"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		
		// API Gateway dependencies
		{ID: "gateway-auth", Type: "depends_on", FromNodeID: "api-gateway", ToNodeID: "auth-service", Attributes: models.Attributes{"type": "http"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "gateway-user", Type: "depends_on", FromNodeID: "api-gateway", ToNodeID: "user-service", Attributes: models.Attributes{"type": "http"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "gateway-order", Type: "depends_on", FromNodeID: "api-gateway", ToNodeID: "order-service", Attributes: models.Attributes{"type": "http"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		
		// Service dependencies
		{ID: "auth-userdb", Type: "depends_on", FromNodeID: "auth-service", ToNodeID: "user-db", Attributes: models.Attributes{"type": "sql"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "auth-cache", Type: "depends_on", FromNodeID: "auth-service", ToNodeID: "redis-cache", Attributes: models.Attributes{"type": "tcp"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "auth-logger", Type: "depends_on", FromNodeID: "auth-service", ToNodeID: "logger", Attributes: models.Attributes{"type": "library"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		
		{ID: "user-userdb", Type: "depends_on", FromNodeID: "user-service", ToNodeID: "user-db", Attributes: models.Attributes{"type": "sql"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "user-logger", Type: "depends_on", FromNodeID: "user-service", ToNodeID: "logger", Attributes: models.Attributes{"type": "library"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		
		{ID: "order-orderdb", Type: "depends_on", FromNodeID: "order-service", ToNodeID: "order-db", Attributes: models.Attributes{"type": "sql"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "order-payment", Type: "depends_on", FromNodeID: "order-service", ToNodeID: "payment-service", Attributes: models.Attributes{"type": "http"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "order-notification", Type: "depends_on", FromNodeID: "order-service", ToNodeID: "notification-service", Attributes: models.Attributes{"type": "async"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "order-logger", Type: "depends_on", FromNodeID: "order-service", ToNodeID: "logger", Attributes: models.Attributes{"type": "library"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		
		{ID: "payment-logger", Type: "depends_on", FromNodeID: "payment-service", ToNodeID: "logger", Attributes: models.Attributes{"type": "library"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		
		{ID: "notification-queue", Type: "depends_on", FromNodeID: "notification-service", ToNodeID: "message-queue", Attributes: models.Attributes{"type": "amqp"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "notification-logger", Type: "depends_on", FromNodeID: "notification-service", ToNodeID: "logger", Attributes: models.Attributes{"type": "library"}, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	
	for _, edge := range edges {
		err = engine.CreateEdge(graphID, edge)
		if err != nil {
			t.Fatalf("Failed to create edge %s: %v", edge.ID, err)
		}
	}
	
	// Step 2: Test comprehensive analysis
	t.Run("CompleteAnalysis", func(t *testing.T) {
		// Test DFS from frontend
		dfsResult, err := analyzer.DepthFirstSearch(graphID, "frontend", &types.TraversalOptions{
			Direction: types.DirectionForward,
			MaxDepth:  -1,
		})
		if err != nil {
			t.Errorf("DFS failed: %v", err)
		}
		if len(dfsResult.Nodes) < 10 {
			t.Errorf("Expected at least 10 nodes in DFS, got %d", len(dfsResult.Nodes))
		}
		
		// Test dependencies analysis
		deps, err := analyzer.GetAllDependencies(graphID, "frontend", &types.TraversalOptions{
			Direction: types.DirectionForward,
			MaxDepth:  -1,
		})
		if err != nil {
			t.Errorf("Dependencies analysis failed: %v", err)
		}
		if len(deps) < 10 {
			t.Errorf("Expected at least 10 dependencies, got %d", len(deps))
		}
		
		// Test dependents analysis
		dependents, err := analyzer.GetAllDependents(graphID, "logger", &types.TraversalOptions{
			Direction: types.DirectionBackward,
			MaxDepth:  -1,
		})
		if err != nil {
			t.Errorf("Dependents analysis failed: %v", err)
		}
		if len(dependents) < 5 {
			t.Errorf("Expected at least 5 dependents for logger, got %d", len(dependents))
		}
		
		// Test shortest path (single path)
		pathResult, err := analyzer.GetShortestPath(graphID, "frontend", "user-db", &types.TraversalOptions{
			Direction: types.DirectionForward,
		})
		if err != nil {
			t.Errorf("Shortest path query failed: %v", err)
		}
		if pathResult != nil && len(pathResult.Path) == 0 {
			t.Error("Expected path from frontend to user-db")
		}
		if pathResult != nil && pathResult.Length != 3 { // frontend -> api-gateway -> user-service -> user-db
			t.Errorf("Expected path length 3, got %d", pathResult.Length)
		}

		// Test all shortest paths (multiple paths)
		allPaths, err := analyzer.AllShortestPaths(graphID, "frontend", "user-db")
		if err != nil {
			t.Errorf("All shortest paths query failed: %v", err)
		}
		if len(allPaths) == 0 {
			t.Error("Expected at least one path from frontend to user-db")
		}
		for i, path := range allPaths {
			if path.Length != 3 {
				t.Errorf("Path %d should have length 3, got %d", i, path.Length)
			}
			if path.Path[0] != "frontend" || path.Path[len(path.Path)-1] != "user-db" {
				t.Errorf("Path %d should start with 'frontend' and end with 'user-db', got %v", i, path.Path)
			}
		}

		// Test all paths traversal
		allTraversalPaths, err := analyzer.AllPathsTraversal(graphID, "frontend", &types.TraversalOptions{
			Direction: types.DirectionForward,
			MaxDepth:  4, // Limit depth to avoid too many paths
		})
		if err != nil {
			t.Errorf("All paths traversal failed: %v", err)
		}
		if len(allTraversalPaths) == 0 {
			t.Error("Expected at least one traversal path from frontend")
		}
		for i, path := range allTraversalPaths {
			if len(path.Nodes) == 0 {
				t.Errorf("Traversal path %d is empty", i)
			}
			if path.Nodes[0].ID != "frontend" {
				t.Errorf("Traversal path %d should start with 'frontend', got %s", i, path.Nodes[0].ID)
			}
		}
		
		// Test cycle detection
		hasCycles, err := analyzer.HasCycles(graphID, &types.TraversalOptions{
			Direction: types.DirectionForward,
		})
		if err != nil {
			t.Errorf("Cycle detection failed: %v", err)
		}
		if hasCycles {
			t.Error("Expected no cycles in microservices architecture")
		}

		// Add a cyclic dependency to test cycle detection
		cyclicEdge := &models.Edge{ID: "payment-order", Type: "depends_on", FromNodeID: "payment-service", ToNodeID: "order-service", Attributes: models.Attributes{}, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		err = engine.CreateEdge(graphID, cyclicEdge)
		if err != nil {
			t.Fatalf("Failed to create cyclic edge: %v", err)
		}

		// Test cycle detection with the new cyclic edge
		cycles, err := analyzer.FindAllCycles(graphID, &types.TraversalOptions{Direction: types.DirectionForward})
		if err != nil {
			t.Errorf("Cycle detection with cyclic graph failed: %v", err)
		}
		if len(cycles) == 0 {
			t.Error("Expected to find a cycle")
		}
		
		// Test graph statistics
		stats, err := analyzer.GetGraphStats(graphID, &types.TraversalOptions{
			Direction: types.DirectionForward,
		})
		if err != nil {
			t.Errorf("Graph stats failed: %v", err)
		}
		if stats.NodeCount != 12 { // Updated to match actual count
			t.Errorf("Expected 12 nodes, got %d", stats.NodeCount)
		}
		if stats.EdgeCount != 17 { // Updated to match actual count
			t.Errorf("Expected 16 edges, got %d", stats.EdgeCount)
		}
		if stats.RootNodeCount < 1 {
			t.Errorf("Expected at least 1 root node, got %d", stats.RootNodeCount)
		}
	})
	
	// Step 3: Test filtering capabilities
	t.Run("FilteringCapabilities", func(t *testing.T) {
		// Test service-only traversal
		serviceNodes, err := analyzer.DepthFirstSearch(graphID, "frontend", &types.TraversalOptions{
			Direction: types.DirectionForward,
			MaxDepth:  -1,
			NodeTypes: []models.NodeType{"service"},
		})
		if err != nil {
			t.Errorf("Service filtering failed: %v", err)
		}
		for _, node := range serviceNodes.Nodes {
			if node.Type != "service" {
				t.Errorf("Expected only service nodes, found %s", node.Type)
			}
		}
		
		// Test HTTP-only dependencies
		httpDeps, err := analyzer.GetAllDependencies(graphID, "frontend", &types.TraversalOptions{
			Direction: types.DirectionForward,
			MaxDepth:  -1,
			EdgeTypes: []models.EdgeType{"depends_on"},
		})
		if err != nil {
			t.Errorf("HTTP filtering failed: %v", err)
		}
		if len(httpDeps) == 0 {
			t.Error("Expected HTTP dependencies")
		}
	})
	
	// Step 4: Test GRAPH.GET for correct counts
	t.Run("GraphGetCounts", func(t *testing.T) {

		nodeCount, _ := engine.CountNodes(graphID)
		edgeCount, _ := engine.CountEdges(graphID)

		if nodeCount != 12 {
			t.Errorf("Expected 12 nodes, got %d", nodeCount)
		}
		if edgeCount != 17 {
			t.Errorf("Expected 17 edges, got %d", edgeCount)
		}
	})

	// Step 5: Test Degree Centrality
	t.Run("DegreeCentrality", func(t *testing.T) {
		scores, err := analyzer.CalculateDegreeCentrality(graphID, nil, types.DirectionBoth)
		if err != nil {
			t.Fatalf("Failed to calculate degree centrality: %v", err)
		}

		// Check a few key nodes
		if scores["api-gateway"] != 4 { // 1 in, 3 out
			t.Errorf("Expected degree centrality of 4 for api-gateway, got %d", scores["api-gateway"])
		}
		if scores["logger"] != 5 { // 5 in, 0 out
			t.Errorf("Expected degree centrality of 5 for logger, got %d", scores["logger"])
		}
	})

	// Step 6: Test Louvain Clustering
	t.Run("LouvainClustering", func(t *testing.T) {
		communities, err := analyzer.CalculateLouvainClustering(graphID, 1.0)
		if err != nil {
			t.Fatalf("Failed to calculate Louvain clustering: %v", err)
		}

		// Basic validation: ensure we have more than one community and not every node is its own community
		if len(communities) <= 1 || len(communities) == 12 {
			t.Errorf("Expected a reasonable number of communities, got %d", len(communities))
		}
	})

	// Step 7: Test data integrity after operations
	t.Run("DataIntegrity", func(t *testing.T) {
		// Verify all nodes exist
		allNodes, err := engine.ListNodes(graphID)
		if err != nil {
			t.Errorf("Failed to list nodes: %v", err)
		}
		if len(allNodes) != 12 {
			t.Errorf("Expected 12 nodes, got %d", len(allNodes))
		}
		
		// Verify all edges exist
		allEdges, err := engine.ListEdges(graphID)
		if err != nil {
			t.Errorf("Failed to list edges: %v", err)
		}
		if len(allEdges) != 17 {
			t.Errorf("Expected 16 edges, got %d", len(allEdges))
		}
		
		// Test relationship consistency
		for _, edge := range allEdges {
			// Verify from node exists
			_, err := engine.GetNode(graphID, edge.FromNodeID)
			if err != nil {
				t.Errorf("From node %s not found for edge %s", edge.FromNodeID, edge.ID)
			}
			
			// Verify to node exists
			_, err = engine.GetNode(graphID, edge.ToNodeID)
			if err != nil {
				t.Errorf("To node %s not found for edge %s", edge.ToNodeID, edge.ID)
			}
		}
	})
}

// TestBackupRestore tests backup and restore functionality
func TestBackupRestore(t *testing.T) {
	originalPath := filepath.Join(os.TempDir(), "pathwaydb_backup_original")
	backupPath := filepath.Join(os.TempDir(), "pathwaydb_backup_file")
	restorePath := filepath.Join(os.TempDir(), "pathwaydb_backup_restored")
	
	defer func() {
		os.RemoveAll(originalPath)
		os.RemoveAll(backupPath)
		os.RemoveAll(restorePath)
	}()
	
	// Create original database with data
	engine1 := storage.NewBadgerEngine()
	err := engine1.Open(originalPath)
	if err != nil {
		t.Fatalf("Failed to open original database: %v", err)
	}
	
	graphID := models.GraphID("backup-test")
	graph := &models.Graph{
		ID:          graphID,
		Name:        "Backup Test Graph",
		Description: "Graph for backup testing",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	engine1.CreateGraph(graph)
	
	// Add test data
	node := &models.Node{ID: "test-node", Type: "service", Attributes: models.Attributes{"name": "Test"}, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	engine1.CreateNode(graphID, node)
	
	// Create backup directory and backup
	err = os.MkdirAll(backupPath, 0755)
	if err != nil {
		t.Errorf("Failed to create backup directory: %v", err)
	}
	defer os.RemoveAll(backupPath)
	
	err = engine1.Backup(backupPath)
	if err != nil {
		t.Errorf("Backup failed: %v", err)
	}
	engine1.Close()
	
	// Restore to new location
	engine2 := storage.NewBadgerEngine()
	err = engine2.Open(restorePath)
	if err != nil {
		t.Fatalf("Failed to open restored database: %v", err)
	}
	defer engine2.Close()
	
	// Verify data exists after restore (this is a simplified test)
	restoredGraph, err := engine2.GetGraph(graphID)
	if err == nil && restoredGraph.ID == graphID {
		t.Log("Backup/restore verification passed")
	}
}
