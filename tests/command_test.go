package tests

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/ywadi/PathwayDB/models"
	"github.com/ywadi/PathwayDB/redis/commands"
	"github.com/ywadi/PathwayDB/storage"
)

// commandTestHarness provides a test environment for command handlers.
 type commandTestHarness struct {
	storage  storage.StorageEngine
	analysis *commands.AnalysisCommands
	testPath string
	graphID  models.GraphID
}

// setupCommandTest creates a new test harness.
func setupCommandTest(t *testing.T) *commandTestHarness {
	testPath := filepath.Join(os.TempDir(), "pathwaydb_command_test_"+t.Name())
	engine := storage.NewBadgerEngine()

	if err := engine.Open(testPath); err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	graphID := models.GraphID("cmd-test-graph")
	graph := &models.Graph{
		ID:   graphID,
		Name: "Command Test Graph",
	}
	if err := engine.CreateGraph(graph); err != nil {
		t.Fatalf("Failed to create graph: %v", err)
	}

	return &commandTestHarness{
		storage:  engine,
		analysis: commands.NewAnalysisCommands(engine),
		testPath: testPath,
		graphID:  graphID,
	}
}

// cleanup closes the engine and removes test data.
func (h *commandTestHarness) cleanup() {
	if h.storage != nil {
		h.storage.Close()
	}
	os.RemoveAll(h.testPath)
}

// createGraph creates a sample graph for testing.
func (h *commandTestHarness) createGraph(nodes []*models.Node, edges []*models.Edge) {
	for _, node := range nodes {
		h.storage.CreateNode(h.graphID, node)
	}
	for _, edge := range edges {
		h.storage.CreateEdge(h.graphID, edge)
	}
}

// TestTraversalCommands tests the ANALYSIS.TRAVERSE command.
func TestTraversalCommands(t *testing.T) {
	h := setupCommandTest(t)
	defer h.cleanup()

	nodes := []*models.Node{
		{ID: "a", Type: "service"},
		{ID: "b", Type: "service"},
		{ID: "c", Type: "database"},
	}
	edges := []*models.Edge{
		{ID: "a-b", FromNodeID: "a", ToNodeID: "b", Type: "calls"},
		{ID: "b-c", FromNodeID: "b", ToNodeID: "c", Type: "writes_to"},
	}
	h.createGraph(nodes, edges)

	t.Run("ForwardTraversal", func(t *testing.T) {
		resp, err := h.analysis.Handle("TRAVERSE", []string{string(h.graphID), "a"})
		if err != nil {
			t.Fatalf("TRAVERSE command failed: %v", err)
		}

		expected := []string{"1", "a:service->a-b:calls->b:service->b-c:writes_to->c:database"}
		if !reflect.DeepEqual(resp.ArrayValue, expected) {
			t.Errorf("Expected response %v, got %v", expected, resp.ArrayValue)
		}
	})

	t.Run("BackwardTraversalWithDirectionalArrows", func(t *testing.T) {
		resp, err := h.analysis.Handle("TRAVERSE", []string{string(h.graphID), "c", "DIRECTION", "in"})
		if err != nil {
			t.Fatalf("TRAVERSE DIRECTION in failed: %v", err)
		}

		expected := []string{"1", "c:database<-b-c:writes_to<-b:service<-a-b:calls<-a:service"}
		if !reflect.DeepEqual(resp.ArrayValue, expected) {
			t.Errorf("Expected response %v, got %v", expected, resp.ArrayValue)
		}
	})

	t.Run("BidirectionalTraversalNoTrivialPaths", func(t *testing.T) {
		resp, err := h.analysis.Handle("TRAVERSE", []string{string(h.graphID), "b", "DIRECTION", "both"})
		if err != nil {
			t.Fatalf("TRAVERSE DIRECTION both failed: %v", err)
		}

		// Expect two paths: b->c and b<-a. No trivial b->a->b path.
		expected := []string{"2", "b:service->b-c:writes_to->c:database", "b:service<-a-b:calls<-a:service"}
		values := resp.ArrayValue
		sort.Strings(values[1:]) // Sort paths for stable comparison
		if !reflect.DeepEqual(values, expected) {
			t.Errorf("Expected response %v, got %v", expected, values)
		}
	})
}

// TestCycleCommands tests the ANALYSIS.CYCLES command.
func TestCycleCommands(t *testing.T) {
	h := setupCommandTest(t)
	defer h.cleanup()

	nodes := []*models.Node{
		{ID: "a", Type: "service"},
		{ID: "b", Type: "service"},
		{ID: "c", Type: "service"},
		{ID: "d", Type: "service"}, // Unrelated to cycles
	}
	edges := []*models.Edge{
		{ID: "a-b", FromNodeID: "a", ToNodeID: "b", Type: "calls"},
		{ID: "b-c", FromNodeID: "b", ToNodeID: "c", Type: "calls"},
		{ID: "c-a", FromNodeID: "c", ToNodeID: "a", Type: "calls"},
		{ID: "c-b", FromNodeID: "c", ToNodeID: "b", Type: "calls"}, // Creates a second, smaller cycle
	}
	h.createGraph(nodes, edges)

	t.Run("FindAllCyclesDetailed", func(t *testing.T) {
		resp, err := h.analysis.Handle("CYCLES", []string{string(h.graphID)})
		if err != nil {
			t.Fatalf("CYCLES command failed: %v", err)
		}

		// Expect 2 cycles: a->b->c->a and b->c->b
		values := resp.ArrayValue
		if values[0] != "2" {
			t.Errorf("Expected to find 2 cycles, got %s", values[0])
		}

		// Normalize and check for expected cycles
		foundC1 := false
		foundC2 := false
		for _, path := range values[1:] {
			if path == "a:service->a-b:calls->b:service->b-c:calls->c:service->c-a:calls->a:service" {
				foundC1 = true
			}
			if path == "b:service->b-c:calls->c:service->c-b:calls->b:service" {
				foundC2 = true
			}
		}

		if !foundC1 || !foundC2 {
			t.Errorf("Did not find all expected cycles. Found C1: %v, Found C2: %v", foundC1, foundC2)
		}
	})

	t.Run("FindAllCyclesSimple", func(t *testing.T) {
		resp, err := h.analysis.Handle("CYCLES", []string{string(h.graphID), "FORMAT", "simple"})
		if err != nil {
			t.Fatalf("CYCLES FORMAT simple failed: %v", err)
		}

		// Expect a, b, c to be in cycles
		expected := []string{"a:service", "b:service", "c:service"}
		values := resp.ArrayValue
		sort.Strings(values)

		if !reflect.DeepEqual(values, expected) {
			t.Errorf("Expected nodes %v, got %v", expected, values)
		}
	})
}
