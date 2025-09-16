package commands

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sort"

	"github.com/ywadi/PathwayDB/analysis"
	"github.com/ywadi/PathwayDB/models"
	"github.com/ywadi/PathwayDB/redis/protocol"
	"github.com/ywadi/PathwayDB/storage"
	"github.com/ywadi/PathwayDB/types"
)

// AnalysisCommands handles analysis-related Redis commands
type AnalysisCommands struct {
	storage  storage.StorageEngine
	analyzer *analysis.GraphAnalyzer
}

// NewAnalysisCommands creates a new analysis commands handler
func NewAnalysisCommands(storageEngine storage.StorageEngine) *AnalysisCommands {
	return &AnalysisCommands{
		storage:  storageEngine,
		analyzer: analysis.NewGraphAnalyzer(storageEngine),
	}
}

// Handle routes analysis commands to their respective handlers
func (a *AnalysisCommands) Handle(command string, args []string) (*protocol.Response, error) {
	switch command {
	case "SHORTESTPATH":
		return a.handleShortestPath(args)
	case "CENTRALITY":
		return a.handleCentrality(args)
	case "CLUSTERING":
		return a.handleClustering(args)
	case "CYCLES":
		return a.handleCycles(args)
	case "TRAVERSE":
		return a.handleTraverse(args)
	default:
		return nil, fmt.Errorf("unknown ANALYSIS command: %s", command)
	}
}

// handleShortestPath handles ANALYSIS.SHORTESTPATH <graph> <from> <to> [algorithm] [FORMAT simple|detailed]
func (a *AnalysisCommands) handleShortestPath(args []string) (*protocol.Response, error) {
	if len(args) < 3 {
		return nil, fmt.Errorf("ANALYSIS.SHORTESTPATH requires at least 3 arguments: graph, from, to")
	}

	graphID := args[0]
	fromNodeID := args[1]
	toNodeID := args[2]

	format := "detailed" // Default to detailed format

	// Parse optional arguments
	for i := 3; i < len(args); i++ {
		if args[i] == "FORMAT" && i+1 < len(args) {
			i++
			if args[i] != "simple" && args[i] != "detailed" {
				return nil, fmt.Errorf("invalid FORMAT: %s (must be 'simple' or 'detailed')", args[i])
			}
			format = args[i]
		}
		// Note: algorithm parameter is parsed but not used yet as the analyzer only supports one algorithm
	}

	// Use the existing GetShortestPath method from GraphAnalyzer
	pathResult, err := a.analyzer.GetShortestPath(models.GraphID(graphID), models.NodeID(fromNodeID), models.NodeID(toNodeID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to compute shortest path: %v", err)
	}

	if pathResult == nil {
		return protocol.NewNullResponse(), nil
	}

	// Simple format with nodeid:nodetype
	if format == "simple" {
		return a.buildSimplePathResponse(models.GraphID(graphID), pathResult)
	}

	// Enhanced detailed format with multiple paths
	allPaths, err := a.analyzer.AllShortestPaths(models.GraphID(graphID), models.NodeID(fromNodeID), models.NodeID(toNodeID))
	if err != nil {
		return nil, fmt.Errorf("failed to find all shortest paths: %v", err)
	}

	if len(allPaths) == 0 {
		return protocol.NewNullResponse(), nil
	}

	return a.buildMultiPathResponse(allPaths)
}

// buildDetailedPathResponse creates a detailed shortest path response with pipe-delimited format
func (a *AnalysisCommands) buildDetailedPathResponse(graphID models.GraphID, pathResult *types.PathResult) (*protocol.Response, error) {
	if len(pathResult.Path) == 0 {
		return protocol.NewNullResponse(), nil
	}

	// Get node details for each node in the path
	nodeDetails := make([]*models.Node, len(pathResult.Path))
	for i, nodeID := range pathResult.Path {
		node, err := a.storage.GetNode(graphID, nodeID)
		if err != nil {
			return nil, fmt.Errorf("failed to get node %s: %v", nodeID, err)
		}
		nodeDetails[i] = node
	}

	// Build arrow notation path: node_id:node_type->edge_id:edge_type->node_id:node_type...
	var pathBuilder strings.Builder

	for i, node := range nodeDetails {
		pathBuilder.WriteString(string(node.ID))
		pathBuilder.WriteString(":")
		pathBuilder.WriteString(string(node.Type))

		if i < len(nodeDetails)-1 {
			// Find edge between current and next node
			nextNode := nodeDetails[i+1]
			edge, err := a.findEdgeBetweenNodes(graphID, node.ID, nextNode.ID)
			if err == nil && edge != nil {
				pathBuilder.WriteString("->")
				pathBuilder.WriteString(string(edge.ID))
				pathBuilder.WriteString(":")
				pathBuilder.WriteString(string(edge.Type))
				pathBuilder.WriteString("->")
			} else {
				pathBuilder.WriteString("->unknown:unknown->")
			}
		}
	}

	// Return single path string
	response := []string{pathBuilder.String()}
	return protocol.NewArrayResponse(response), nil
}

// handleCentrality handles ANALYSIS.CENTRALITY <graph> <type> [node_id] [DIRECTION in|out|both]
// type can be: "betweenness", "closeness", "degree"
func (a *AnalysisCommands) handleCentrality(args []string) (*protocol.Response, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("ANALYSIS.CENTRALITY requires at least 2 arguments: graph, type")
	}

	graphID := models.GraphID(args[0])
	centralityType := strings.ToLower(args[1])

	var nodeID *models.NodeID
	direction := types.DirectionBoth // Default direction

	// Parse optional arguments: node_id and DIRECTION
	i := 2
	for i < len(args) {
		if strings.ToUpper(args[i]) == "DIRECTION" {
			if i+1 >= len(args) {
				return nil, fmt.Errorf("DIRECTION option requires an argument")
			}
			i++
			switch strings.ToLower(args[i]) {
			case "in":
				direction = types.DirectionBackward
			case "out":
				direction = types.DirectionForward
			case "both":
				direction = types.DirectionBoth
			default:
				return nil, fmt.Errorf("invalid DIRECTION: %s", args[i])
			}
			i++
		} else {
			if nodeID != nil {
				return nil, fmt.Errorf("unexpected argument: %s. node_id already provided", args[i])
			}
			tempNodeID := models.NodeID(args[i])
			nodeID = &tempNodeID
			i++
		}
	}

	switch centralityType {
	case "degree":
		scores, err := a.analyzer.CalculateDegreeCentrality(graphID, nodeID, direction)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate degree centrality: %w", err)
		}

		response := make([]string, 0, len(scores)*2)
		for id, score := range scores {
			response = append(response, string(id), strconv.Itoa(score))
		}
		return protocol.NewArrayResponse(response), nil
	case "betweenness", "closeness":
		// TODO: Implement betweenness and closeness centrality
		return protocol.NewArrayResponse([]string{"centrality", centralityType, "not_implemented"}), nil
	default:
		return nil, fmt.Errorf("unknown centrality type: %s", centralityType)
	}
}

// handleClustering handles ANALYSIS.CLUSTERING <graph> [algorithm] [parameters_json]
func (a *AnalysisCommands) handleClustering(args []string) (*protocol.Response, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("ANALYSIS.CLUSTERING requires at least 1 argument: graph")
	}

	graphID := models.GraphID(args[0])
	algorithm := "louvain" // default
	if len(args) > 1 {
		algorithm = args[1]
	}

	// Default parameters
	resolution := 1.0

	// Parse parameters if provided
	if len(args) > 2 {
		var params map[string]interface{}
		if err := json.Unmarshal([]byte(args[2]), &params); err != nil {
			return nil, fmt.Errorf("invalid parameters JSON: %v", err)
		}
		if res, ok := params["resolution"]; ok {
			if r, ok := res.(float64); ok {
				resolution = r
			} else {
				return nil, fmt.Errorf("resolution parameter must be a float")
			}
		}
	}

	switch algorithm {
	case "louvain":
		communities, err := a.analyzer.CalculateLouvainClustering(graphID, resolution)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate Louvain clustering: %w", err)
		}

		// Format response as an array of communities, where each community is an array of node IDs.
		// This requires a custom response builder as it's a nested array.
		communityArrays := make([]interface{}, len(communities))
		for i, community := range communities {
			nodeIDs := make([]string, len(community))
			for j, nodeID := range community {
				nodeIDs[j] = string(nodeID)
			}
			communityArrays[i] = nodeIDs
		}
		return protocol.NewNestedArrayResponse(communityArrays), nil

	case "connected_components":
		componentCount, err := a.analyzer.GetConnectedComponentCount(models.GraphID(graphID), nil)
		if err != nil {
			return nil, fmt.Errorf("failed to compute connected components: %v", err)
		}
		result := []string{"connected_components", strconv.Itoa(componentCount)}
		return protocol.NewArrayResponse(result), nil

	default:
		return nil, fmt.Errorf("unknown clustering algorithm: %s", algorithm)
	}
}

// handleCycles handles ANALYSIS.CYCLES <graph> [NODETYPE type1...] [EDGETYPE type1...] [FORMAT simple|detailed]
func (a *AnalysisCommands) handleCycles(args []string) (*protocol.Response, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("ANALYSIS.CYCLES requires at least 1 argument: graph")
	}

	graphID := args[0]
	format := "detailed" // Default to detailed format
	options := &types.TraversalOptions{
		Direction: types.DirectionForward,
	}

	// Parse optional filters
	i := 1
	for i < len(args) {
		switch args[i] {
		case "NODETYPE", "NODETYPES":
			i++
			for i < len(args) && args[i] != "EDGETYPE" && args[i] != "EDGETYPES" && args[i] != "FORMAT" {
				options.NodeTypes = append(options.NodeTypes, models.NodeType(args[i]))
				i++
			}
		case "EDGETYPE", "EDGETYPES":
			i++
			for i < len(args) && args[i] != "NODETYPE" && args[i] != "NODETYPES" && args[i] != "FORMAT" {
				options.EdgeTypes = append(options.EdgeTypes, models.EdgeType(args[i]))
				i++
			}
		case "FORMAT":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("FORMAT option requires an argument")
			}
			i++
			if args[i] != "simple" && args[i] != "detailed" {
				return nil, fmt.Errorf("invalid FORMAT: %s (must be 'simple' or 'detailed')", args[i])
			}
			format = args[i]
			i++
		default:
			return nil, fmt.Errorf("unknown option for ANALYSIS.CYCLES: %s", args[i])
		}
	}

	cycles, err := a.analyzer.FindAllCycles(models.GraphID(graphID), options)
	if err != nil {
		return nil, fmt.Errorf("failed to check for cycles: %v", err)
	}

	if len(cycles) == 0 {
		return protocol.NewNullResponse(), nil
	}

	// Return response based on format
	if format == "simple" {
		uniqueNodes := make(map[models.NodeID]bool)
		for _, cycle := range cycles {
			// The last node is a repeat of the first, so we can skip it.
			for _, nodeID := range cycle[:len(cycle)-1] {
				uniqueNodes[nodeID] = true
			}
		}

		response := make([]string, 0, len(uniqueNodes))
		for nodeID := range uniqueNodes {
			node, err := a.storage.GetNode(models.GraphID(graphID), nodeID)
			if err != nil {
				return nil, fmt.Errorf("failed to get node %s: %v", nodeID, err)
			}
			response = append(response, string(node.ID)+":"+string(node.Type))
		}

		// Sort for deterministic output
		sort.Strings(response)

		return protocol.NewArrayResponse(response), nil
	}

	return a.buildDetailedCycleResponse(models.GraphID(graphID), cycles)
}

// buildSimpleCycleResponse creates a simple cycle response with nodeid:nodetype format
func (a *AnalysisCommands) buildSimpleCycleResponse(graphID models.GraphID, cyclePath []models.NodeID) (*protocol.Response, error) {
	if len(cyclePath) == 0 {
		return protocol.NewNullResponse(), nil
	}

	response := make([]string, len(cyclePath))
	for i, nodeID := range cyclePath {
		node, err := a.storage.GetNode(graphID, nodeID)
		if err != nil {
			return nil, fmt.Errorf("failed to get node %s: %v", nodeID, err)
		}
		response[i] = string(node.ID) + ":" + string(node.Type)
	}

	return protocol.NewArrayResponse(response), nil
}

// buildDetailedCycleResponse creates a detailed response for multiple cycles with arrow notation.
func (a *AnalysisCommands) buildDetailedCycleResponse(graphID models.GraphID, cycles [][]models.NodeID) (*protocol.Response, error) {
	if len(cycles) == 0 {
		return protocol.NewNullResponse(), nil
	}

	var cycleStrings []string
	for _, cyclePath := range cycles {
		if len(cyclePath) == 0 {
			continue
		}

		nodeDetails := make([]*models.Node, len(cyclePath))
		for i, nodeID := range cyclePath {
			node, err := a.storage.GetNode(graphID, nodeID)
			if err != nil {
				return nil, fmt.Errorf("failed to get node %s: %v", nodeID, err)
			}
			nodeDetails[i] = node
		}

		var pathBuilder strings.Builder
		for i := 0; i < len(nodeDetails)-1; i++ {
			currentNode := nodeDetails[i]
			nextNode := nodeDetails[i+1]

			pathBuilder.WriteString(string(currentNode.ID))
			pathBuilder.WriteString(":")
			pathBuilder.WriteString(string(currentNode.Type))

			edge, err := a.findEdgeBetweenNodes(graphID, currentNode.ID, nextNode.ID)
			if err == nil && edge != nil {
				arrow := buildArrow(currentNode.ID, nextNode.ID, edge)
				pathBuilder.WriteString(arrow)
				pathBuilder.WriteString(string(edge.ID))
				pathBuilder.WriteString(":")
				pathBuilder.WriteString(string(edge.Type))
				pathBuilder.WriteString(arrow)
			} else {
				pathBuilder.WriteString("->unknown:unknown->")
			}
		}

		lastNode := nodeDetails[len(nodeDetails)-1]
		pathBuilder.WriteString(string(lastNode.ID))
		pathBuilder.WriteString(":")
		pathBuilder.WriteString(string(lastNode.Type))

		cycleStrings = append(cycleStrings, pathBuilder.String())
	}

	// Return cycle count and the detailed path strings
	response := make([]string, 0, len(cycleStrings)+1)
	response = append(response, strconv.Itoa(len(cycleStrings)))
	response = append(response, cycleStrings...)

	return protocol.NewArrayResponse(response), nil
}

// handleTraverse handles ANALYSIS.TRAVERSE <graph> <start_node> [DIRECTION dir] [NODETYPES type1...] [EDGETYPES type1...] [FORMAT simple|detailed]
func (a *AnalysisCommands) handleTraverse(args []string) (*protocol.Response, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("ANALYSIS.TRAVERSE requires at least 2 arguments: graph, start_node")
	}

	graphID := models.GraphID(args[0])
	startNodeID := models.NodeID(args[1])

	options := &types.TraversalOptions{
		Direction: types.DirectionForward, // Default direction
		MaxDepth:  -1,                     // No depth limit for full traversal
	}

	format := "detailed" // Default to detailed format

	// Parse optional keyword arguments
	i := 2
	for i < len(args) {
		switch args[i] {
		case "DIRECTION":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("DIRECTION option requires an argument")
			}
			i++
			switch args[i] {
			case "in":
				options.Direction = types.DirectionBackward
			case "out":
				options.Direction = types.DirectionForward
			case "both":
				options.Direction = types.DirectionBoth
			default:
				return nil, fmt.Errorf("invalid DIRECTION: %s", args[i])
			}
			i++
		case "NODETYPES":
			i++
			// Accept multiple node types (OR logic)
			for i < len(args) && args[i] != "EDGETYPES" && args[i] != "DIRECTION" && args[i] != "FORMAT" {
				options.NodeTypes = append(options.NodeTypes, models.NodeType(args[i]))
				i++
			}
		case "EDGETYPES":
			i++
			// Accept multiple edge types (OR logic)
			for i < len(args) && args[i] != "NODETYPES" && args[i] != "DIRECTION" && args[i] != "FORMAT" {
				options.EdgeTypes = append(options.EdgeTypes, models.EdgeType(args[i]))
				i++
			}
		case "FORMAT":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("FORMAT option requires an argument")
			}
			i++
			if args[i] != "simple" && args[i] != "detailed" {
				return nil, fmt.Errorf("invalid FORMAT: %s (must be 'simple' or 'detailed')", args[i])
			}
			format = args[i]
			i++
		default:
			return nil, fmt.Errorf("unknown option for ANALYSIS.TRAVERSE: %s", args[i])
		}
	}

	// Use AllPathsTraversal for detailed format to get multiple paths
	if format == "detailed" {
		allPaths, err := a.analyzer.AllPathsTraversal(models.GraphID(graphID), startNodeID, options)
		if err != nil {
			return nil, fmt.Errorf("failed to perform multi-path traversal: %v", err)
		}

		if allPaths == nil || len(allPaths) == 0 {
			return protocol.NewNullResponse(), nil
		}

		return a.buildMultiPathTraversalResponse(allPaths)
	}

	// Use single path traversal for simple format
	result, err := a.analyzer.DepthFirstSearch(models.GraphID(graphID), startNodeID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to traverse graph: %v", err)
	}

	if result == nil {
		return protocol.NewNullResponse(), nil
	}

	return a.buildSimpleTraversalResponse(result)
}

// buildSimpleTraversalResponse creates a simple traversal response with nodeid:nodetype format
func (a *AnalysisCommands) buildSimpleTraversalResponse(result *types.TraversalResult) (*protocol.Response, error) {
	if len(result.Nodes) == 0 {
		return protocol.NewNullResponse(), nil
	}

	response := make([]string, len(result.Nodes))
	for i, node := range result.Nodes {
		response[i] = string(node.ID) + ":" + string(node.Type)
	}

	return protocol.NewArrayResponse(response), nil
}

// buildSimplePathResponse creates a simple path response with nodeid:nodetype format
func (a *AnalysisCommands) buildSimplePathResponse(graphID models.GraphID, pathResult *types.PathResult) (*protocol.Response, error) {
	if len(pathResult.Path) == 0 {
		return protocol.NewNullResponse(), nil
	}

	response := make([]string, len(pathResult.Path))
	for i, nodeID := range pathResult.Path {
		// Get node details
		node, err := a.storage.GetNode(graphID, nodeID)
		if err != nil {
			return nil, fmt.Errorf("failed to get node %s: %v", nodeID, err)
		}

		response[i] = string(node.ID) + ":" + string(node.Type)
	}

	return protocol.NewArrayResponse(response), nil
}

// buildDetailedTraversalResponse creates a detailed traversal response with pipe-delimited format
func (a *AnalysisCommands) buildDetailedTraversalResponse(graphID models.GraphID, result *types.TraversalResult) (*protocol.Response, error) {
	if len(result.Nodes) == 0 {
		return protocol.NewNullResponse(), nil
	}

	// For traversal, we need to handle multiple paths from the starting node
	// Group paths by their starting edges to handle branching
	paths := make([]string, 0)

	// If we have a linear path (single traversal path)
	if len(result.Nodes) > 0 {
		var pathBuilder strings.Builder

		for i, node := range result.Nodes {
			pathBuilder.WriteString(string(node.ID))
			pathBuilder.WriteString(":")
			pathBuilder.WriteString(string(node.Type))

			if i < len(result.Nodes)-1 {
				// Find edge to next node
				var edge *models.Edge
				if i < len(result.Edges) {
					edge = result.Edges[i]
				} else {
					// Fallback to finding edge between consecutive nodes
					nextNode := result.Nodes[i+1]
					edge, _ = a.findEdgeBetweenNodes(graphID, node.ID, nextNode.ID)
				}

				if edge != nil {
					pathBuilder.WriteString("->")
					pathBuilder.WriteString(string(edge.ID))
					pathBuilder.WriteString(":")
					pathBuilder.WriteString(string(edge.Type))
					pathBuilder.WriteString("->")
				} else {
					pathBuilder.WriteString("->unknown:unknown->")
				}
			}
		}

		paths = append(paths, pathBuilder.String())
	}

	// Return paths count followed by each path
	response := make([]string, 0, len(paths)+1)
	response = append(response, fmt.Sprintf("%d", len(paths)))
	response = append(response, paths...)

	return protocol.NewArrayResponse(response), nil
}

// buildArrow determines the correct arrow notation based on traversal direction.
func buildArrow(fromNode, toNode models.NodeID, edge *models.Edge) string {
	if edge.FromNodeID == fromNode && edge.ToNodeID == toNode {
		return "->"
	} else if edge.ToNodeID == fromNode && edge.FromNodeID == toNode {
		return "<-"
	}
	return "->" // Default for safety, though this case should be rare.
}

// buildMultiPathTraversalResponse creates response for multiple traversal paths
func (a *AnalysisCommands) buildMultiPathTraversalResponse(allPaths []*types.TraversalResult) (*protocol.Response, error) {
	response := make([]string, 0, len(allPaths)+1)
	response = append(response, fmt.Sprintf("%d", len(allPaths)))

	for _, path := range allPaths {
		var pathBuilder strings.Builder

		for i, node := range path.Nodes {
			pathBuilder.WriteString(string(node.ID))
			pathBuilder.WriteString(":")
			pathBuilder.WriteString(string(node.Type))

			if i < len(path.Nodes)-1 && i < len(path.Edges) {
				edge := path.Edges[i]
				arrow := buildArrow(node.ID, path.Nodes[i+1].ID, edge)

				pathBuilder.WriteString(arrow)
				pathBuilder.WriteString(string(edge.ID))
				pathBuilder.WriteString(":")
				pathBuilder.WriteString(string(edge.Type))
				pathBuilder.WriteString(arrow)
			}
		}

		response = append(response, pathBuilder.String())
	}

	return protocol.NewArrayResponse(response), nil
}

// buildMultiPathResponse creates response for multiple shortest paths
func (a *AnalysisCommands) buildMultiPathResponse(allPaths []*types.PathResult) (*protocol.Response, error) {
	response := make([]string, 0, len(allPaths)+1)
	response = append(response, fmt.Sprintf("%d", len(allPaths)))

	for _, pathResult := range allPaths {
		// Get node details for each node in the path
		nodeDetails := make([]*models.Node, len(pathResult.Path))
		for i, nodeID := range pathResult.Path {
			node, err := a.storage.GetNode(models.GraphID(pathResult.FromNodeID), nodeID)
			if err != nil {
				return nil, fmt.Errorf("failed to get node %s: %v", nodeID, err)
			}
			nodeDetails[i] = node
		}

		// Build arrow notation path
		var pathBuilder strings.Builder

		for i, node := range nodeDetails {
			pathBuilder.WriteString(string(node.ID))
			pathBuilder.WriteString(":")
			pathBuilder.WriteString(string(node.Type))

			if i < len(nodeDetails)-1 && i < len(pathResult.Edges) {
				// Get edge details
				edgeID := pathResult.Edges[i]
				edge, err := a.storage.GetEdge(models.GraphID(pathResult.FromNodeID), edgeID)
				if err == nil && edge != nil {
					arrow := buildArrow(node.ID, nodeDetails[i+1].ID, edge)
					pathBuilder.WriteString(arrow)
					pathBuilder.WriteString(string(edge.ID))
					pathBuilder.WriteString(":")
					pathBuilder.WriteString(string(edge.Type))
					pathBuilder.WriteString(arrow)
				} else {
					pathBuilder.WriteString("->unknown:unknown->")
				}
			}
		}

		response = append(response, pathBuilder.String())
	}

	return protocol.NewArrayResponse(response), nil
}

// findEdgeBetweenNodes finds an edge between two nodes
func (a *AnalysisCommands) findEdgeBetweenNodes(graphID models.GraphID, fromNode, toNode models.NodeID) (*models.Edge, error) {
	// Get outgoing edges from the source node
	outgoingEdges, err := a.storage.GetOutgoingEdges(graphID, fromNode)
	if err != nil {
		return nil, fmt.Errorf("failed to get outgoing edges from %s: %v", fromNode, err)
	}

	// Find edge that connects to the target node
	for _, edge := range outgoingEdges {
		if edge.ToNodeID == toNode {
			return edge, nil
		}
	}

	// Also check incoming edges to the target node (which would be outgoing from other nodes)
	incomingEdges, err := a.storage.GetIncomingEdges(graphID, toNode)
	if err != nil {
		return nil, fmt.Errorf("failed to get incoming edges to %s: %v", toNode, err)
	}

	for _, edge := range incomingEdges {
		if edge.FromNodeID == fromNode {
			return edge, nil
		}
	}

	return nil, fmt.Errorf("no edge found between nodes %s and %s", fromNode, toNode)
}
