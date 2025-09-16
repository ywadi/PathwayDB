package analysis

import (
	"fmt"

	"github.com/ywadi/PathwayDB/models"
	"github.com/ywadi/PathwayDB/storage"
	"github.com/ywadi/PathwayDB/types"
	"gonum.org/v1/gonum/graph/community"
	"strings"

	"gonum.org/v1/gonum/graph/simple"
)

// GraphAnalyzer provides comprehensive graph analysis capabilities
type GraphAnalyzer struct {
	storage storage.StorageEngine
}

// NewGraphAnalyzer creates a new graph analyzer instance
func NewGraphAnalyzer(storage storage.StorageEngine) *GraphAnalyzer {
	return &GraphAnalyzer{
		storage: storage,
	}
}

// DepthFirstSearch performs a depth-first search traversal starting from a given node
func (ga *GraphAnalyzer) DepthFirstSearch(graphID models.GraphID, startNodeID models.NodeID, options *types.TraversalOptions) (*types.TraversalResult, error) {
	if options == nil {
		options = &types.TraversalOptions{
			MaxDepth:  -1, // No limit
			Direction: types.DirectionForward,
		}
	}

	visited := make(map[models.NodeID]bool)
	var nodes []*models.Node
	var edges []*models.Edge
	var path []models.NodeID

	// Use iterative DFS with a stack
	type stackItem struct {
		nodeID models.NodeID
		depth  int
	}

	stack := []stackItem{{nodeID: startNodeID, depth: 0}}

	for len(stack) > 0 {
		// Main traversal loop
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// Skip if already visited or depth exceeded
		if visited[current.nodeID] {
			continue
		}
		if options.MaxDepth >= 0 && current.depth > options.MaxDepth {
			continue
		}

		// Mark as visited
		visited[current.nodeID] = true

		// Get the current node
		node, err := ga.storage.GetNode(graphID, current.nodeID)
		if err != nil {
			return nil, fmt.Errorf("failed to get node %s: %w", current.nodeID, err)
		}

		// Check stop condition
		if options.StopCondition != nil && options.StopCondition(node) {
			continue
		}

		// Check node type filter
		nodeTypeMatch := true
		if len(options.NodeTypes) > 0 {
			nodeTypeMatch = false
			for _, nodeType := range options.NodeTypes {
				if node.Type == nodeType {
					nodeTypeMatch = true
					break
				}
			}
		}

		// Add to results if node type matches
		if nodeTypeMatch {
			nodes = append(nodes, node)
			path = append(path, current.nodeID)
		}

		// Get connected edges
		var connectedEdges []*models.Edge
		switch options.Direction {
		case types.DirectionForward:
			connectedEdges, err = ga.storage.GetOutgoingEdges(graphID, current.nodeID)
		case types.DirectionBackward:
			connectedEdges, err = ga.storage.GetIncomingEdges(graphID, current.nodeID)
		case types.DirectionBoth:
			outgoing, err1 := ga.storage.GetOutgoingEdges(graphID, current.nodeID)
			if err1 != nil {
				return nil, fmt.Errorf("failed to get outgoing edges: %w", err1)
			}
			incoming, err2 := ga.storage.GetIncomingEdges(graphID, current.nodeID)
			if err2 != nil {
				return nil, fmt.Errorf("failed to get incoming edges: %w", err2)
			}
			connectedEdges = append(outgoing, incoming...)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to get connected edges: %w", err)
		}

		// Filter edges by type if specified
		if len(options.EdgeTypes) > 0 {
			var filteredEdges []*models.Edge
			for _, edge := range connectedEdges {
				for _, edgeType := range options.EdgeTypes {
					if edge.Type == edgeType {
						filteredEdges = append(filteredEdges, edge)
						break
					}
				}
			}
			connectedEdges = filteredEdges
		}

		// Add edges and connected nodes to stack (in reverse order for DFS)
		for i := len(connectedEdges) - 1; i >= 0; i-- {
			edge := connectedEdges[i]

			// Determine next node
			var nextNodeID models.NodeID
			switch options.Direction {
			case types.DirectionForward:
				if edge.FromNodeID == current.nodeID {
					nextNodeID = edge.ToNodeID
				}
			case types.DirectionBackward:
				if edge.ToNodeID == current.nodeID {
					nextNodeID = edge.FromNodeID
				}
			case types.DirectionBoth:
				if edge.FromNodeID == current.nodeID {
					nextNodeID = edge.ToNodeID
				} else if edge.ToNodeID == current.nodeID {
					nextNodeID = edge.FromNodeID
				}
			}

			// Add to stack if not visited
			if nextNodeID != "" && !visited[nextNodeID] {
				edges = append(edges, edge)
				stack = append(stack, stackItem{
					nodeID: nextNodeID,
					depth:  current.depth + 1,
				})
			}
		}
	}

	return &types.TraversalResult{
		Nodes:    nodes,
		Edges:    edges,
		Path:     path,
		Distance: len(path) - 1,
	}, nil
}

// AllPathsTraversal finds all complete paths from a starting node, exploring all branches
func (ga *GraphAnalyzer) AllPathsTraversal(graphID models.GraphID, startNodeID models.NodeID, options *types.TraversalOptions) ([]*types.TraversalResult, error) {
	if options == nil {
		options = &types.TraversalOptions{
			MaxDepth:  -1, // No limit
			Direction: types.DirectionForward,
		}
	}

	var allPaths []*types.TraversalResult
	visited := make(map[models.NodeID]bool)

	// Start recursive path finding
	err := ga.findAllPathsRecursive(graphID, startNodeID, "", visited, []models.NodeID{}, []*models.Edge{}, 0, options, &allPaths)
	if err != nil {
		return nil, err
	}

	return allPaths, nil
}

// findAllPathsRecursive recursively finds all paths from current node
func (ga *GraphAnalyzer) findAllPathsRecursive(graphID models.GraphID, nodeID models.NodeID, previousEdgeID models.EdgeID, visited map[models.NodeID]bool,
	currentPath []models.NodeID, currentEdges []*models.Edge, depth int, options *types.TraversalOptions, allPaths *[]*types.TraversalResult) error {

	// Check depth limit
	if options.MaxDepth >= 0 && depth > options.MaxDepth {
		return nil
	}

	// Skip if already visited in this path (prevent cycles)
	if visited[nodeID] {
		return nil
	}

	// Get the current node
	node, err := ga.storage.GetNode(graphID, nodeID)
	if err != nil {
		return fmt.Errorf("failed to get node %s: %w", nodeID, err)
	}

	// Check stop condition
	if options.StopCondition != nil && options.StopCondition(node) {
		return nil
	}

	// Check node type filter
	nodeTypeMatch := true
	if len(options.NodeTypes) > 0 {
		nodeTypeMatch = false
		for _, nodeType := range options.NodeTypes {
			if node.Type == nodeType {
				nodeTypeMatch = true
				break
			}
		}
	}

	// Add current node to path if it matches filter
	if nodeTypeMatch {
		currentPath = append(currentPath, nodeID)
	}

	// Mark as visited for this path
	visited[nodeID] = true
	defer func() { visited[nodeID] = false }() // Unmark when backtracking

	// Get connected edges
	var connectedEdges []*models.Edge
	switch options.Direction {
	case types.DirectionForward:
		connectedEdges, err = ga.storage.GetOutgoingEdges(graphID, nodeID)
	case types.DirectionBackward:
		connectedEdges, err = ga.storage.GetIncomingEdges(graphID, nodeID)
	case types.DirectionBoth:
		outgoing, err1 := ga.storage.GetOutgoingEdges(graphID, nodeID)
		if err1 != nil {
			return fmt.Errorf("failed to get outgoing edges: %w", err1)
		}
		incoming, err2 := ga.storage.GetIncomingEdges(graphID, nodeID)
		if err2 != nil {
			return fmt.Errorf("failed to get incoming edges: %w", err2)
		}
		connectedEdges = append(outgoing, incoming...)
	}

	if err != nil {
		return fmt.Errorf("failed to get connected edges: %w", err)
	}

	// Filter edges by type if specified
	if len(options.EdgeTypes) > 0 {
		var filteredEdges []*models.Edge
		for _, edge := range connectedEdges {
			for _, edgeType := range options.EdgeTypes {
				if edge.Type == edgeType {
					filteredEdges = append(filteredEdges, edge)
					break
				}
			}
		}
		connectedEdges = filteredEdges
	}


	// In 'both' direction, we need to filter out the edge we just came from
	// before deciding if this is a leaf node.
	if options.Direction == types.DirectionBoth && previousEdgeID != "" {
		var potentialNextEdges []*models.Edge
		for _, edge := range connectedEdges {
			if edge.ID != previousEdgeID {
				potentialNextEdges = append(potentialNextEdges, edge)
			}
		}
		connectedEdges = potentialNextEdges
	}

	// If no outgoing edges, this is a leaf node - save the current path
	if len(connectedEdges) == 0 {
		if len(currentPath) > 0 {
			// Convert path to nodes
			pathNodes := make([]*models.Node, len(currentPath))
			for i, pathNodeID := range currentPath {
				pathNode, err := ga.storage.GetNode(graphID, pathNodeID)
				if err != nil {
					return fmt.Errorf("failed to get path node %s: %w", pathNodeID, err)
				}
				pathNodes[i] = pathNode
			}

			*allPaths = append(*allPaths, &types.TraversalResult{
				Nodes:    pathNodes,
				Edges:    append([]*models.Edge{}, currentEdges...), // Copy edges
				Path:     append([]models.NodeID{}, currentPath...), // Copy path
				Distance: len(currentPath) - 1,
			})
		}
		return nil
	}

	// Explore each connected edge
	for _, edge := range connectedEdges {
		// Determine next node
		var nextNodeID models.NodeID
		switch options.Direction {
		case types.DirectionForward:
			if edge.FromNodeID == nodeID {
				nextNodeID = edge.ToNodeID
			}
		case types.DirectionBackward:
			if edge.ToNodeID == nodeID {
				nextNodeID = edge.FromNodeID
			}
		case types.DirectionBoth:
			if edge.FromNodeID == nodeID {
				nextNodeID = edge.ToNodeID
			} else if edge.ToNodeID == nodeID {
				nextNodeID = edge.FromNodeID
			}
		}

		// Continue recursion if next node is valid
		if nextNodeID != "" {
			newEdges := append(currentEdges, edge)

			// If the neighbor is already in the path, we have a cycle.
			if visited[nextNodeID] {
				cycleStartIndex := -1
				for i, pathNodeID := range currentPath {
					if pathNodeID == nextNodeID {
						cycleStartIndex = i
						break
					}
				}

				if cycleStartIndex != -1 {
					// Construct the cycle path and edges
					cyclePath := currentPath[cycleStartIndex:]
					cycleEdges := newEdges[cycleStartIndex:]

					// Get node objects for the path
					pathNodes := make([]*models.Node, len(cyclePath))
					for i, pathNodeID := range cyclePath {
						pathNode, nodeErr := ga.storage.GetNode(graphID, pathNodeID)
						if nodeErr != nil {
							return fmt.Errorf("failed to get cycle path node %s: %w", pathNodeID, nodeErr)
						}
						pathNodes[i] = pathNode
					}

					// Add the closing node to complete the cycle visualization
					pathNodes = append(pathNodes, pathNodes[0])
					cyclePath = append(cyclePath, cyclePath[0])

					*allPaths = append(*allPaths, &types.TraversalResult{
						Nodes:    pathNodes,
						Edges:    cycleEdges,
						Path:     cyclePath,
						Distance: len(cyclePath) - 1,
					})
				}
			} else {
				// Continue recursion if it's not a cycle
				err := ga.findAllPathsRecursive(graphID, nextNodeID, edge.ID, visited, currentPath, newEdges, depth+1, options, allPaths)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// dfsRecursive performs the recursive DFS traversal
func (ga *GraphAnalyzer) dfsRecursive(graphID models.GraphID, nodeID models.NodeID, visited map[models.NodeID]bool, nodes *[]*models.Node, edges *[]*models.Edge, path *[]models.NodeID, depth int, options *types.TraversalOptions) error {
	// Check depth limit
	if options.MaxDepth >= 0 && depth > options.MaxDepth {
		return nil
	}

	// Skip if already visited
	if visited[nodeID] {
		return nil
	}

	// Mark as visited immediately to prevent cycles
	visited[nodeID] = true

	// Get the current node
	node, err := ga.storage.GetNode(graphID, nodeID)
	if err != nil {
		return fmt.Errorf("failed to get node %s: %w", nodeID, err)
	}

	// Check stop condition
	if options.StopCondition != nil && options.StopCondition(node) {
		return nil
	}

	// Check node type filter and add to results if it matches
	nodeTypeMatch := true
	if len(options.NodeTypes) > 0 {
		nodeTypeMatch = false
		for _, nodeType := range options.NodeTypes {
			if node.Type == nodeType {
				nodeTypeMatch = true
				break
			}
		}
	}

	// Add to results only if node type matches (or no filter specified)
	if nodeTypeMatch {
		*nodes = append(*nodes, node)
		*path = append(*path, nodeID)
	}

	// Get connected edges based on direction
	var connectedEdges []*models.Edge
	switch options.Direction {
	case types.DirectionForward:
		connectedEdges, err = ga.storage.GetOutgoingEdges(graphID, nodeID)
		if err != nil {
			return fmt.Errorf("failed to get outgoing edges: %w", err)
		}
	case types.DirectionBackward:
		connectedEdges, err = ga.storage.GetIncomingEdges(graphID, nodeID)
		if err != nil {
			return fmt.Errorf("failed to get incoming edges: %w", err)
		}
	case types.DirectionBoth:
		outgoing, err1 := ga.storage.GetOutgoingEdges(graphID, nodeID)
		if err1 != nil {
			return fmt.Errorf("failed to get outgoing edges: %w", err1)
		}
		incoming, err2 := ga.storage.GetIncomingEdges(graphID, nodeID)
		if err2 != nil {
			return fmt.Errorf("failed to get incoming edges: %w", err2)
		}
		connectedEdges = append(outgoing, incoming...)
	}

	// Filter edges by type if specified
	if len(options.EdgeTypes) > 0 {
		var filteredEdges []*models.Edge
		for _, edge := range connectedEdges {
			for _, edgeType := range options.EdgeTypes {
				if edge.Type == edgeType {
					filteredEdges = append(filteredEdges, edge)
					break
				}
			}
		}
		connectedEdges = filteredEdges
	}

	// Traverse connected nodes
	for _, edge := range connectedEdges {
		// Determine next node based on direction and current node
		var nextNodeID models.NodeID
		switch options.Direction {
		case types.DirectionForward:
			if edge.FromNodeID == nodeID {
				nextNodeID = edge.ToNodeID
			}
		case types.DirectionBackward:
			if edge.ToNodeID == nodeID {
				nextNodeID = edge.FromNodeID
			}
		case types.DirectionBoth:
			if edge.FromNodeID == nodeID {
				nextNodeID = edge.ToNodeID
			} else if edge.ToNodeID == nodeID {
				nextNodeID = edge.FromNodeID
			}
		}

		// Recursively traverse the next node
		if nextNodeID != "" && !visited[nextNodeID] {
			// Add edge to results before traversing
			*edges = append(*edges, edge)
			err = ga.dfsRecursive(graphID, nextNodeID, visited, nodes, edges, path, depth+1, options)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetAllDependencies returns a flat list of all transitive dependencies
func (ga *GraphAnalyzer) GetAllDependencies(graphID models.GraphID, nodeID models.NodeID, options *types.TraversalOptions) ([]*models.Node, error) {
	if options == nil {
		options = &types.TraversalOptions{
			Direction: types.DirectionForward,
		}
	}

	result, err := ga.DepthFirstSearch(graphID, nodeID, options)
	if err != nil {
		return nil, err
	}

	// Remove the starting node from dependencies
	var dependencies []*models.Node
	for _, node := range result.Nodes {
		if node.ID != nodeID {
			dependencies = append(dependencies, node)
		}
	}

	return dependencies, nil
}

// GetAllDependents returns a flat list of all transitive dependents
func (ga *GraphAnalyzer) GetAllDependents(graphID models.GraphID, nodeID models.NodeID, options *types.TraversalOptions) ([]*models.Node, error) {
	if options == nil {
		options = &types.TraversalOptions{
			Direction: types.DirectionBackward,
		}
	} else {
		// Override direction for dependents
		options.Direction = types.DirectionBackward
	}

	result, err := ga.DepthFirstSearch(graphID, nodeID, options)
	if err != nil {
		return nil, err
	}

	// Remove the starting node from dependents
	var dependents []*models.Node
	for _, node := range result.Nodes {
		if node.ID != nodeID {
			dependents = append(dependents, node)
		}
	}

	return dependents, nil
}

// GetShortestPath finds the shortest path between two nodes using BFS
func (ga *GraphAnalyzer) GetShortestPath(graphID models.GraphID, fromNodeID, toNodeID models.NodeID, options *types.TraversalOptions) (*types.PathResult, error) {
	if options == nil {
		options = &types.TraversalOptions{
			Direction: types.DirectionForward,
		}
	}

	visited := make(map[models.NodeID]bool)
	parent := make(map[models.NodeID]models.NodeID)
	edgeMap := make(map[models.NodeID]models.EdgeID)

	// Queue for BFS
	type queueItem struct {
		nodeID models.NodeID
		depth  int
	}
	queue := []queueItem{{nodeID: fromNodeID, depth: 0}}
	visited[fromNodeID] = true

	found := false
	for len(queue) > 0 && !found {
		current := queue[0]
		queue = queue[1:]

		if current.nodeID == toNodeID {
			found = true
			break
		}

		// Get connected edges
		var connectedEdges []*models.Edge
		var err error
		switch options.Direction {
		case types.DirectionForward:
			connectedEdges, err = ga.storage.GetOutgoingEdges(graphID, current.nodeID)
		case types.DirectionBackward:
			connectedEdges, err = ga.storage.GetIncomingEdges(graphID, current.nodeID)
		case types.DirectionBoth:
			outgoing, err1 := ga.storage.GetOutgoingEdges(graphID, current.nodeID)
			if err1 != nil {
				return nil, fmt.Errorf("failed to get outgoing edges: %w", err1)
			}
			incoming, err2 := ga.storage.GetIncomingEdges(graphID, current.nodeID)
			if err2 != nil {
				return nil, fmt.Errorf("failed to get incoming edges: %w", err2)
			}
			connectedEdges = append(outgoing, incoming...)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to get connected edges: %w", err)
		}

		for _, edge := range connectedEdges {
			var nextNodeID models.NodeID
			switch options.Direction {
			case types.DirectionForward:
				if edge.FromNodeID == current.nodeID {
					nextNodeID = edge.ToNodeID
				}
			case types.DirectionBackward:
				if edge.ToNodeID == current.nodeID {
					nextNodeID = edge.FromNodeID
				}
			case types.DirectionBoth:
				if edge.FromNodeID == current.nodeID {
					nextNodeID = edge.ToNodeID
				} else if edge.ToNodeID == current.nodeID {
					nextNodeID = edge.FromNodeID
				}
			}

			if nextNodeID != "" && !visited[nextNodeID] {
				visited[nextNodeID] = true
				parent[nextNodeID] = current.nodeID
				edgeMap[nextNodeID] = edge.ID
				queue = append(queue, queueItem{
					nodeID: nextNodeID,
					depth:  current.depth + 1,
				})
			}
		}
	}

	if !found {
		return nil, fmt.Errorf("no path found from %s to %s", fromNodeID, toNodeID)
	}

	// Reconstruct path
	var path []models.NodeID
	var edges []models.EdgeID
	currentNode := toNodeID

	for currentNode != fromNodeID {
		path = append([]models.NodeID{currentNode}, path...)
		if edgeID, exists := edgeMap[currentNode]; exists {
			edges = append([]models.EdgeID{edgeID}, edges...)
		}
		currentNode = parent[currentNode]
	}
	path = append([]models.NodeID{fromNodeID}, path...)

	return &types.PathResult{
		FromNodeID: fromNodeID,
		ToNodeID:   toNodeID,
		Path:       path,
		Length:     len(path) - 1,
		Edges:      edges,
	}, nil
}

// AllShortestPaths finds all shortest paths between two nodes
func (ga *GraphAnalyzer) AllShortestPaths(graphID models.GraphID, fromNodeID, toNodeID models.NodeID) ([]*types.PathResult, error) {
	// Use BFS to find all paths of minimum length
	type queueItem struct {
		nodeID models.NodeID
		path   []models.NodeID
		edges  []*models.Edge
		dist   int
	}

	queue := []queueItem{{nodeID: fromNodeID, path: []models.NodeID{fromNodeID}, edges: []*models.Edge{}, dist: 0}}
	visited := make(map[models.NodeID]int) // Track minimum distance to each node
	var allPaths []*types.PathResult
	minDistance := -1

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// If we've found a path and current distance is greater, stop
		if minDistance != -1 && current.dist > minDistance {
			break
		}

		// If we reached the target
		if current.nodeID == toNodeID {
			if minDistance == -1 {
				minDistance = current.dist
			}

			// Only add paths of minimum distance
			if current.dist == minDistance {
				// Convert edges to edge IDs
				edgeIDs := make([]models.EdgeID, len(current.edges))
				for i, edge := range current.edges {
					edgeIDs[i] = edge.ID
				}

				allPaths = append(allPaths, &types.PathResult{
					FromNodeID: fromNodeID,
					ToNodeID:   toNodeID,
					Path:       append([]models.NodeID{}, current.path...),
					Length:     current.dist,
					Edges:      edgeIDs,
				})
			}
			continue
		}

		// Skip if we've seen this node with a shorter distance
		if prevDist, seen := visited[current.nodeID]; seen && prevDist < current.dist {
			continue
		}
		visited[current.nodeID] = current.dist

		// Get outgoing edges
		outgoingEdges, err := ga.storage.GetOutgoingEdges(graphID, current.nodeID)
		if err != nil {
			return nil, fmt.Errorf("failed to get outgoing edges from %s: %w", current.nodeID, err)
		}

		// Explore neighbors
		for _, edge := range outgoingEdges {
			nextNodeID := edge.ToNodeID

			// Avoid cycles in the current path
			inCurrentPath := false
			for _, pathNode := range current.path {
				if pathNode == nextNodeID {
					inCurrentPath = true
					break
				}
			}

			if !inCurrentPath {
				newPath := append(current.path, nextNodeID)
				newEdges := append(current.edges, edge)
				queue = append(queue, queueItem{
					nodeID: nextNodeID,
					path:   newPath,
					edges:  newEdges,
					dist:   current.dist + 1,
				})
			}
		}
	}

	return allPaths, nil
}

// FindAllCycles finds all elementary cycles in the graph.
func (ga *GraphAnalyzer) FindAllCycles(graphID models.GraphID, options *types.TraversalOptions) ([][]models.NodeID, error) {
	if options == nil {
		options = &types.TraversalOptions{
			Direction: types.DirectionForward,
		}
	}

	allNodes, err := ga.storage.ListNodes(graphID)
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes for cycle detection: %w", err)
	}

	var allCycles [][]models.NodeID
	for _, node := range allNodes {
		path := []models.NodeID{node.ID}
		blocked := make(map[models.NodeID]bool)
		cycles, err := ga.findCyclesRecursive(graphID, node.ID, node.ID, path, blocked, &allCycles, options)
		if err != nil {
			return nil, err
		}
		allCycles = append(allCycles, cycles...)
	}

	// Deduplicate cycles
	uniqueCycles := make(map[string][]models.NodeID)
	for _, cycle := range allCycles {
		normalized := normalizeCycle(cycle)
		uniqueCycles[cycleToString(normalized)] = normalized
	}

	result := make([][]models.NodeID, 0, len(uniqueCycles))
	for _, cycle := range uniqueCycles {
		result = append(result, cycle)
	}

	return result, nil
}

// normalizeCycle creates a canonical representation of a cycle by rotating it
// to start with the lexicographically smallest node ID.
func normalizeCycle(path []models.NodeID) []models.NodeID {
	if len(path) <= 1 {
		return path
	}

	// The actual cycle nodes, excluding the repeated start node at the end
	cycleNodes := path[:len(path)-1]

	// Find the index of the lexicographically smallest node
	minIndex := 0
	for i := 1; i < len(cycleNodes); i++ {
		if cycleNodes[i] < cycleNodes[minIndex] {
			minIndex = i
		}
	}

	// Rotate the slice to start with the smallest node
	rotated := append(cycleNodes[minIndex:], cycleNodes[:minIndex]...)
	// Append the start node to close the cycle
	return append(rotated, rotated[0])
}

// cycleToString converts a cycle path to a unique string identifier.
func cycleToString(path []models.NodeID) string {
	ids := make([]string, len(path))
	for i, nodeID := range path {
		ids[i] = string(nodeID)
	}
	return strings.Join(ids, "->")
}

func (ga *GraphAnalyzer) findCyclesRecursive(graphID models.GraphID, startNode, currentNode models.NodeID, path []models.NodeID, blocked map[models.NodeID]bool, allCycles *[][]models.NodeID, options *types.TraversalOptions) ([][]models.NodeID, error) {
	var newCycles [][]models.NodeID
	blocked[currentNode] = true
	defer func() { blocked[currentNode] = false }() // Unblock node on backtrack

	connectedEdges, err := ga.storage.GetOutgoingEdges(graphID, currentNode)
	if err != nil {
		return nil, fmt.Errorf("failed to get outgoing edges from %s: %w", currentNode, err)
	}

	// Filter edges by type if specified
	if options != nil && len(options.EdgeTypes) > 0 {
		var filteredEdges []*models.Edge
		for _, edge := range connectedEdges {
			for _, edgeType := range options.EdgeTypes {
				if edge.Type == edgeType {
					filteredEdges = append(filteredEdges, edge)
					break
				}
			}
		}
		connectedEdges = filteredEdges
	}

	for _, edge := range connectedEdges {
		neighbor := edge.ToNodeID
		if neighbor == startNode {
			cycle := make([]models.NodeID, len(path), len(path)+1)
			copy(cycle, path)
			cycle = append(cycle, startNode)
			*allCycles = append(*allCycles, cycle)
			newCycles = append(newCycles, cycle)
		} else if !blocked[neighbor] {
			newPath := append(path, neighbor)
			cycles, err := ga.findCyclesRecursive(graphID, startNode, neighbor, newPath, blocked, allCycles, options)
			if err != nil {
				return nil, err
			}
			newCycles = append(newCycles, cycles...)
		}
	}


	return newCycles, nil
}

// HasCycles checks if the graph contains any cycles by calling FindAllCycles.
func (ga *GraphAnalyzer) HasCycles(graphID models.GraphID, options *types.TraversalOptions) (bool, error) {
	cycles, err := ga.FindAllCycles(graphID, options)
	if err != nil {
		return false, err
	}
	return len(cycles) > 0, nil
}

// dfsHasCycle performs DFS to detect cycles using the three-color approach
func (ga *GraphAnalyzer) dfsHasCycle(graphID models.GraphID, nodeID models.NodeID, visited map[models.NodeID]int, options *types.TraversalOptions) (bool, error) {
	// Mark as visiting (gray)
	visited[nodeID] = 1

	// Get outgoing edges
	var connectedEdges []*models.Edge
	var err error

	switch options.Direction {
	case types.DirectionForward:
		connectedEdges, err = ga.storage.GetOutgoingEdges(graphID, nodeID)
	case types.DirectionBackward:
		connectedEdges, err = ga.storage.GetIncomingEdges(graphID, nodeID)
	case types.DirectionBoth:
		outgoing, err1 := ga.storage.GetOutgoingEdges(graphID, nodeID)
		if err1 != nil {
			return false, fmt.Errorf("failed to get outgoing edges: %w", err1)
		}
		incoming, err2 := ga.storage.GetIncomingEdges(graphID, nodeID)
		if err2 != nil {
			return false, fmt.Errorf("failed to get incoming edges: %w", err2)
		}
		connectedEdges = append(outgoing, incoming...)
	}

	if err != nil {
		return false, fmt.Errorf("failed to get connected edges: %w", err)
	}

	// Filter edges by type if specified
	if len(options.EdgeTypes) > 0 {
		var filteredEdges []*models.Edge
		for _, edge := range connectedEdges {
			for _, edgeType := range options.EdgeTypes {
				if edge.Type == edgeType {
					filteredEdges = append(filteredEdges, edge)
					break
				}
			}
		}
		connectedEdges = filteredEdges
	}

	// Check each connected node
	for _, edge := range connectedEdges {
		var nextNodeID models.NodeID
		switch options.Direction {
		case types.DirectionForward:
			if edge.FromNodeID == nodeID {
				nextNodeID = edge.ToNodeID
			}
		case types.DirectionBackward:
			if edge.ToNodeID == nodeID {
				nextNodeID = edge.FromNodeID
			}
		case types.DirectionBoth:
			if edge.FromNodeID == nodeID {
				nextNodeID = edge.ToNodeID
			} else if edge.ToNodeID == nodeID {
				nextNodeID = edge.FromNodeID
			}
		}

		if nextNodeID != "" {
			if visited[nextNodeID] == 1 {
				// Found a back edge - cycle detected
				return true, nil
			}
			if visited[nextNodeID] == 0 {
				hasCycle, err := ga.dfsHasCycle(graphID, nextNodeID, visited, options)
				if err != nil {
					return false, err
				}
				if hasCycle {
					return true, nil
				}
			}
		}
	}

	// Mark as visited (black)
	visited[nodeID] = 2
	return false, nil
}

// GetGraphStats calculates comprehensive statistics for a graph
func (ga *GraphAnalyzer) GetGraphStats(graphID models.GraphID, options *types.TraversalOptions) (*types.GraphStats, error) {
	// Get all nodes and edges
	allNodes, err := ga.storage.ListNodes(graphID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	allEdges, err := ga.storage.ListEdges(graphID)
	if err != nil {
		return nil, fmt.Errorf("failed to get edges: %w", err)
	}

	stats := &types.GraphStats{
		NodeCount:     len(allNodes),
		EdgeCount:     len(allEdges),
		NodeTypeCount: make(map[models.NodeType]int),
		EdgeTypeCount: make(map[models.EdgeType]int),
	}

	// Count node types
	for _, node := range allNodes {
		stats.NodeTypeCount[node.Type]++
	}

	// Count edge types
	for _, edge := range allEdges {
		stats.EdgeTypeCount[edge.Type]++
	}

	// Calculate root nodes (nodes with no incoming edges)
	rootNodes, err := ga.GetRootNodes(graphID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get root nodes: %w", err)
	}
	stats.RootNodeCount = len(rootNodes)

	// Calculate leaf nodes (nodes with no outgoing edges)
	leafNodes, err := ga.GetLeafNodes(graphID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get leaf nodes: %w", err)
	}
	stats.LeafNodeCount = len(leafNodes)

	// Calculate orphan nodes (nodes with no connections)
	orphanNodes, err := ga.GetOrphanNodes(graphID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get orphan nodes: %w", err)
	}
	stats.OrphanNodeCount = len(orphanNodes)

	// Check for cycles
	hasCycles, err := ga.HasCycles(graphID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to check for cycles: %w", err)
	}
	stats.HasCycles = hasCycles

	// Calculate maximum depth
	maxDepth, err := ga.GetMaxDepth(graphID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate max depth: %w", err)
	}
	stats.MaxDepth = maxDepth

	// Calculate connected components
	componentCount, err := ga.GetConnectedComponentCount(graphID, options)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate connected components: %w", err)
	}
	stats.ConnectedComponents = componentCount

	return stats, nil
}

// GetRootNodes returns nodes with no incoming edges (dependencies)
func (ga *GraphAnalyzer) GetRootNodes(graphID models.GraphID, options *types.TraversalOptions) ([]*models.Node, error) {
	allNodes, err := ga.storage.ListNodes(graphID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	var rootNodes []*models.Node
	for _, node := range allNodes {
		incomingEdges, err := ga.storage.GetIncomingEdges(graphID, node.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get incoming edges for node %s: %w", node.ID, err)
		}

		// Filter edges by type if specified
		if options != nil && len(options.EdgeTypes) > 0 {
			var filteredEdges []*models.Edge
			for _, edge := range incomingEdges {
				for _, edgeType := range options.EdgeTypes {
					if edge.Type == edgeType {
						filteredEdges = append(filteredEdges, edge)
						break
					}
				}
			}
			incomingEdges = filteredEdges
		}

		if len(incomingEdges) == 0 {
			rootNodes = append(rootNodes, node)
		}
	}

	return rootNodes, nil
}

// GetLeafNodes returns nodes with no outgoing edges (dependents)
func (ga *GraphAnalyzer) GetLeafNodes(graphID models.GraphID, options *types.TraversalOptions) ([]*models.Node, error) {
	allNodes, err := ga.storage.ListNodes(graphID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	var leafNodes []*models.Node
	for _, node := range allNodes {
		outgoingEdges, err := ga.storage.GetOutgoingEdges(graphID, node.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get outgoing edges for node %s: %w", node.ID, err)
		}

		// Filter edges by type if specified
		if options != nil && len(options.EdgeTypes) > 0 {
			var filteredEdges []*models.Edge
			for _, edge := range outgoingEdges {
				for _, edgeType := range options.EdgeTypes {
					if edge.Type == edgeType {
						filteredEdges = append(filteredEdges, edge)
						break
					}
				}
			}
			outgoingEdges = filteredEdges
		}

		if len(outgoingEdges) == 0 {
			leafNodes = append(leafNodes, node)
		}
	}

	return leafNodes, nil
}

// GetOrphanNodes returns nodes with no connections (neither incoming nor outgoing edges)
// CalculateDegreeCentrality calculates the degree centrality for nodes in the graph.
// If a specific nodeID is provided, it calculates for that node only.
// Otherwise, it calculates for all nodes in the graph.
// Direction can be 'in', 'out', or 'both'.
func (ga *GraphAnalyzer) CalculateDegreeCentrality(graphID models.GraphID, nodeID *models.NodeID, direction types.TraversalDirection) (map[models.NodeID]int, error) {
	scores := make(map[models.NodeID]int)

	nodesToProcess := []*models.Node{}
	if nodeID != nil {
		node, err := ga.storage.GetNode(graphID, *nodeID)
		if err != nil {
			return nil, fmt.Errorf("failed to get node %s: %w", *nodeID, err)
		}
		nodesToProcess = append(nodesToProcess, node)
	} else {
		var err error
		nodesToProcess, err = ga.storage.ListNodes(graphID)
		if err != nil {
			return nil, fmt.Errorf("failed to list nodes: %w", err)
		}
	}

	for _, node := range nodesToProcess {
		var degree int
		if direction == types.DirectionForward || direction == types.DirectionBoth {
			outgoing, err := ga.storage.GetOutgoingEdges(graphID, node.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get outgoing edges for %s: %w", node.ID, err)
			}
			degree += len(outgoing)
		}
		if direction == types.DirectionBackward || direction == types.DirectionBoth {
			incoming, err := ga.storage.GetIncomingEdges(graphID, node.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get incoming edges for %s: %w", node.ID, err)
			}
			degree += len(incoming)
		}
		scores[node.ID] = degree
	}

	return scores, nil
}

// CalculateLouvainClustering performs community detection using the Louvain algorithm.
func (ga *GraphAnalyzer) CalculateLouvainClustering(graphID models.GraphID, resolution float64) ([][]models.NodeID, error) {
	gonumGraph, nodeMap, err := ga.toGonumGraph(graphID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to gonum graph: %w", err)
	}

	// The community.Modularize function performs the Louvain community detection.
	// It returns a ReducedGraph, which represents the graph with communities as nodes.
	communitiesResult := community.Modularize(gonumGraph, resolution, nil)

	// The Communities method on the result gives us the list of communities.
	gonumCommunities := communitiesResult.Communities()

	// Convert the gonum node IDs back to our internal model's NodeID.
	communities := make([][]models.NodeID, len(gonumCommunities))
	for i, c := range gonumCommunities {
		communityNodes := make([]models.NodeID, len(c))
		for j, node := range c {
			for modelID, gonumID := range nodeMap {
				if gonumID == node.ID() {
					communityNodes[j] = modelID
					break
				}
			}
		}
		communities[i] = communityNodes
	}

	return communities, nil
}

// toGonumGraph converts our internal graph representation to a gonum graph.
func (ga *GraphAnalyzer) toGonumGraph(graphID models.GraphID) (*simple.UndirectedGraph, map[models.NodeID]int64, error) {
	nodes, err := ga.storage.ListNodes(graphID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	edges, err := ga.storage.ListEdges(graphID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list edges: %w", err)
	}

	gonumGraph := simple.NewUndirectedGraph()
	nodeMap := make(map[models.NodeID]int64)
	var nextID int64

	// Add nodes to the graph and create a mapping from our NodeID to gonum's int64 ID.
	for _, node := range nodes {
		if _, exists := nodeMap[node.ID]; !exists {
			nodeMap[node.ID] = nextID
			gonumGraph.AddNode(simple.Node(nextID))
			nextID++
		}
	}

	// Add edges to the graph.
	for _, edge := range edges {
		fromID, fromExists := nodeMap[edge.FromNodeID]
		toID, toExists := nodeMap[edge.ToNodeID]

		if fromExists && toExists {
			gonumGraph.SetEdge(simple.Edge{F: simple.Node(fromID), T: simple.Node(toID)})
		}
	}

	return gonumGraph, nodeMap, nil
}

func (ga *GraphAnalyzer) GetOrphanNodes(graphID models.GraphID, options *types.TraversalOptions) ([]*models.Node, error) {
	allNodes, err := ga.storage.ListNodes(graphID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	var orphanNodes []*models.Node
	for _, node := range allNodes {
		incomingEdges, err := ga.storage.GetIncomingEdges(graphID, node.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get incoming edges for node %s: %w", node.ID, err)
		}

		outgoingEdges, err := ga.storage.GetOutgoingEdges(graphID, node.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get outgoing edges for node %s: %w", node.ID, err)
		}

		if len(incomingEdges) == 0 && len(outgoingEdges) == 0 {
			orphanNodes = append(orphanNodes, node)
		}
	}

	return orphanNodes, nil
}

// GetMaxDepth calculates the maximum depth of the dependency tree
func (ga *GraphAnalyzer) GetMaxDepth(graphID models.GraphID, options *types.TraversalOptions) (int, error) {
	rootNodes, err := ga.GetRootNodes(graphID, options)
	if err != nil {
		return 0, fmt.Errorf("failed to get root nodes: %w", err)
	}

	maxDepth := 0
	for _, rootNode := range rootNodes {
		depth, err := ga.calculateNodeDepth(graphID, rootNode.ID, make(map[models.NodeID]bool), 0, options)
		if err != nil {
			return 0, fmt.Errorf("failed to calculate depth for root node %s: %w", rootNode.ID, err)
		}
		if depth > maxDepth {
			maxDepth = depth
		}
	}

	return maxDepth, nil
}

// calculateNodeDepth recursively calculates the depth of a node
func (ga *GraphAnalyzer) calculateNodeDepth(graphID models.GraphID, nodeID models.NodeID, visited map[models.NodeID]bool, currentDepth int, options *types.TraversalOptions) (int, error) {
	if visited[nodeID] {
		return currentDepth, nil
	}

	visited[nodeID] = true
	maxChildDepth := currentDepth

	outgoingEdges, err := ga.storage.GetOutgoingEdges(graphID, nodeID)
	if err != nil {
		return 0, fmt.Errorf("failed to get outgoing edges: %w", err)
	}

	for _, edge := range outgoingEdges {
		childDepth, err := ga.calculateNodeDepth(graphID, edge.ToNodeID, visited, currentDepth+1, options)
		if err != nil {
			return 0, err
		}
		if childDepth > maxChildDepth {
			maxChildDepth = childDepth
		}
	}

	visited[nodeID] = false
	return maxChildDepth, nil
}

// GetConnectedComponentCount calculates the number of connected components in the graph
func (ga *GraphAnalyzer) GetConnectedComponentCount(graphID models.GraphID, options *types.TraversalOptions) (int, error) {
	allNodes, err := ga.storage.ListNodes(graphID)
	if err != nil {
		return 0, fmt.Errorf("failed to get nodes: %w", err)
	}

	visited := make(map[models.NodeID]bool)
	componentCount := 0

	for _, node := range allNodes {
		if !visited[node.ID] {
			componentCount++
			err := ga.markConnectedComponent(graphID, node.ID, visited, options)
			if err != nil {
				return 0, fmt.Errorf("failed to mark connected component: %w", err)
			}
		}
	}

	return componentCount, nil
}

// markConnectedComponent marks all nodes in a connected component as visited
func (ga *GraphAnalyzer) markConnectedComponent(graphID models.GraphID, nodeID models.NodeID, visited map[models.NodeID]bool, options *types.TraversalOptions) error {
	if visited[nodeID] {
		return nil
	}

	visited[nodeID] = true

	// Get all connected edges (both incoming and outgoing)
	outgoingEdges, err := ga.storage.GetOutgoingEdges(graphID, nodeID)
	if err != nil {
		return fmt.Errorf("failed to get outgoing edges: %w", err)
	}

	incomingEdges, err := ga.storage.GetIncomingEdges(graphID, nodeID)
	if err != nil {
		return fmt.Errorf("failed to get incoming edges: %w", err)
	}

	allEdges := append(outgoingEdges, incomingEdges...)

	// Recursively mark connected nodes
	for _, edge := range allEdges {
		if edge.FromNodeID == nodeID && !visited[edge.ToNodeID] {
			err := ga.markConnectedComponent(graphID, edge.ToNodeID, visited, options)
			if err != nil {
				return err
			}
		}
		if edge.ToNodeID == nodeID && !visited[edge.FromNodeID] {
			err := ga.markConnectedComponent(graphID, edge.FromNodeID, visited, options)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
