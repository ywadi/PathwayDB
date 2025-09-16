# PathwayDB Test Suite

This directory contains comprehensive tests for the PathwayDB graph database system.

## Test Structure

### `storage_test.go`
Tests the storage layer (BadgerEngine) functionality:
- **Graph CRUD**: Create, Read, Update, Delete operations for graphs
- **Node CRUD**: Complete node lifecycle management with attributes and type filtering
- **Edge CRUD**: Edge operations with relationship queries and attribute filtering
- **Database Operations**: Open, close, backup, transaction management
- **Error Handling**: Invalid operations, non-existent resources, closed database scenarios
- **TTL**: Node and edge expiration, including cascading deletes for nodes.

### `analysis_test.go`
Tests the analysis engine functionality:
- **Depth-First Search**: Basic DFS, depth limits, filtering by node/edge types, directional traversal
- **Dependency Analysis**: Transitive dependencies and dependents with filtering
- **Shortest Path**: Path finding, non-existent paths, same-node scenarios
- **Cycle Detection**: Acyclic graphs, cyclic graphs, self-loops
- **Graph Statistics**: Node counts, edge counts, root/leaf/orphan nodes, connected components
- **Node Classification**: Root, leaf, and orphan node identification with filtering
- **Graph Metrics**: Max depth calculation, connected component counting
- **Error Handling**: Empty graphs, non-existent nodes, nil options

### `integration_test.go`
End-to-end integration tests:
- **Complete Workflow**: Full microservices architecture simulation
- **Complex Analysis**: Multi-level dependency analysis on realistic data
- **Filtering Capabilities**: Advanced filtering across node and edge types
- **Data Integrity**: Consistency verification after complex operations
- **Backup/Restore**: Database backup and restoration functionality

## Test Coverage

The test suite covers:

### Storage Layer Functions (BadgerEngine)
- ✅ CreateGraph, GetGraph, UpdateGraph, DeleteGraph, ListGraphs
- ✅ CreateNode, GetNode, UpdateNode, DeleteNode, ListNodes, ListNodesByType, FindNodesByAttribute
- ✅ CreateEdge, GetEdge, UpdateEdge, DeleteEdge, ListEdges, ListEdgesByType
- ✅ GetOutgoingEdges, GetIncomingEdges, GetConnectedNodes, FindEdgesByAttribute
- ✅ Open, Close, Backup, RunTransaction, RunReadOnlyTransaction
- ✅ TTL expiration for nodes and edges

### Analysis Engine Functions (GraphAnalyzer)
- ✅ DepthFirstSearch with all options (direction, depth, filtering)
- ✅ GetAllDependencies, GetAllDependents with filtering
- ✅ GetShortestPath with various scenarios
- ✅ HasCycles for acyclic and cyclic graphs
- ✅ GetGraphStats with comprehensive metrics
- ✅ GetRootNodes, GetLeafNodes, GetOrphanNodes with filtering
- ✅ GetMaxDepth, GetConnectedComponentCount

### Edge Cases and Error Scenarios
- ✅ Operations on closed databases
- ✅ Non-existent graphs, nodes, and edges
- ✅ Invalid edge creation (non-existent nodes)
- ✅ Duplicate resource creation
- ✅ Empty graphs and disconnected components
- ✅ Nil options handling
- ✅ Self-loops and cycles
- ✅ Transaction rollback scenarios

## Running Tests

```bash
# Run all tests
go test ./tests/... -v

# Run specific test files
go test ./tests/storage_test.go -v
go test ./tests/analysis_test.go -v
go test ./tests/integration_test.go -v

# Run with coverage
go test ./tests/... -cover -v

# Run specific test functions
go test ./tests/ -run TestGraphCRUD -v
go test ./tests/ -run TestDepthFirstSearch -v
```

## Test Data

Tests use temporary databases that are automatically cleaned up after each test. The test suite includes:

- **Simple graphs**: Basic A->B->C chains for algorithm verification
- **Complex dependency graphs**: Realistic microservices architectures
- **Cyclic graphs**: For cycle detection testing
- **Disconnected components**: For connectivity analysis
- **Various node/edge types**: For filtering and classification testing

## Assertions and Validations

Each test includes comprehensive assertions for:
- **Correctness**: Results match expected values
- **Completeness**: All expected data is returned
- **Consistency**: Data integrity is maintained
- **Error Handling**: Appropriate errors for invalid operations
- **Performance**: Operations complete within reasonable time
- **Memory**: No memory leaks in long-running tests
