# PathwayDB - A High-Performance Graph & Dependency Database

[![Docker Build and Test](https://github.com/ywadi/PathwayDB/actions/workflows/docker-build.yml/badge.svg)](https://github.com/ywadi/PathwayDB/actions/workflows/docker-build.yml)

PathwayDB is a high-performance, embeddable graph database built in Go. It uses Badger v3 as its underlying storage engine, providing a fast, transactional, and persistent graph data store. It is designed for modeling complex dependencies, such as microservices architecture, infrastructure, or any other system with interconnected components.

## Features

- **High-Performance**: Built on Badger v3's LSM-tree for fast reads and writes.
- **Embeddable & Pure Go**: No external dependencies or CGO, easy to integrate into any Go application.
- **ACID Transactions**: Guarantees atomicity, consistency, isolation, and durability for all operations.
- **Rich Data Model**: Supports typed nodes and edges with flexible key-value attributes.
- **Powerful Analysis Engine**: Includes advanced graph traversal and analysis algorithms.
- **Redis-Compatible Protocol**: Interact with the database using a namespaced Redis-compatible API.
- **Web-Based IDE**: A modern, professional IDE for real-time graph visualization and command execution.
- **Time-To-Live (TTL) with Cascading Deletes**: Set an expiration on nodes and edges. Expired nodes will be automatically deleted along with their connected edges.
- **Comprehensive Test Suite**: Ensures reliability and correctness with high test coverage.

## Project Structure

```
PathwayDB/
├── README.md           # Project documentation
├── go.mod              # Go module definition
├── go.sum              # Dependency checksums
├── analysis/           # Graph analysis engine
├── cmd/                # Server executables
│   ├── ide-server/
│   └── redis-server/
├── data/               # Default data directory for BadgerDB
├── ide/                # Web-based IDE (React frontend, Go backend)
├── models/             # Core data models (Graph, Node, Edge)
├── redis/              # Redis protocol implementation
├── storage/            # Storage engine implementation
├── tests/              # Comprehensive test suite
├── types/              # Analysis-related type definitions
└── utils/              # Key encoding and other utilities
```

## Getting Started

There are three main ways to use PathwayDB:

1.  **As an embedded Go library**
2.  **As a standalone Redis-compatible server**
3.  **Through the Web IDE**

### 1. Using the IDE (Recommended)

The easiest way to get started is with the integrated web IDE. It provides a comprehensive graphical interface for interacting with the database.

```bash
# From the project root, run the start script
./ide/start.sh
```

This script will:
1.  Start the PathwayDB Redis server.
2.  Start the IDE's WebSocket backend server.
3.  Build and start the React frontend.
4.  Load sample data into the database.

Once running, you can access the IDE at **http://localhost:3000**.

### 2. Running the Standalone Redis Server

You can run PathwayDB as a standalone server that speaks the Redis protocol.

```bash
# Build the server
go build -o redis-server ./cmd/redis-server

# Run the server
./redis-server
```

By default, the server listens on port `6379`. You can connect to it using any standard Redis client, such as `redis-cli`.

### 3. Using as a Go Library

To use PathwayDB in your own Go project, simply import the `storage` and `analysis` packages.

```go
package main

import (
	"fmt"
	"github.com/ywadi/PathwayDB/analysis"
	"github.com/ywadi/PathwayDB/storage"
)

func main() {
	// Initialize the storage engine
	db := storage.NewBadgerEngine()
	err := db.Open("/tmp/pathwaydb_docs")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Initialize the analysis engine
	analyzer := analysis.NewGraphAnalyzer(db)

	// Use the analyzer to perform graph operations...
	fmt.Println("PathwayDB is ready!")
}
```

## PathwayDB IDE

The IDE provides a modern, professional interface for managing and visualizing your graphs.

- **Professional UI**: Modern dark theme with compact, business-ready design.
- **Graph Visualization**: Interactive graph visualization using Cytoscape.js with multiple layouts.
- **Redis Console**: Real-time command execution with syntax highlighting and command history.
- **WebSocket Integration**: Direct connection to the Redis protocol for real-time updates.
- **Graph Explorer**: Browse graphs, nodes, and edges with detailed statistics.
- **Properties Panel**: Inspect selected nodes and edges with full attribute details and TTL information.

## Redis Protocol Reference

All PathwayDB commands are namespaced to avoid conflicts with standard Redis commands. The available namespaces are `GRAPH`, `NODE`, `EDGE`, and `ANALYSIS`.

### `GRAPH` Commands

- `GRAPH.CREATE <name> [description]`
- `GRAPH.DELETE <name>`
- `GRAPH.LIST`
- `GRAPH.GET <name>`
- `GRAPH.EXISTS <name>`

### `NODE` Commands

- `NODE.CREATE <graph> <id> <type> [attributes_json] [TTL <seconds>]`
- `NODE.GET <graph> <id>`
- `NODE.UPDATE <graph> <id> <attributes_json> [TTL <seconds>]`
- `NODE.DELETE <graph> <id>`
- `NODE.FILTER <graph> <attribute_key> <attribute_value>`
- `NODE.LIST <graph>`
- `NODE.EXISTS <graph> <id>`

### `EDGE` Commands

- `EDGE.CREATE <graph> <id> <from> <to> <type> [attributes_json] [TTL <seconds>]`
- `EDGE.GET <graph> <id>`
- `EDGE.UPDATE <graph> <id> <attributes_json> [TTL <seconds>]`
- `EDGE.DELETE <graph> <id>`
- `EDGE.FILTER <graph> <attribute_key> <attribute_value>`
- `EDGE.NEIGHBORS <graph> <node> [DIRECTION in|out|both] [FORMAT simple|detailed]`
- `EDGE.LIST <graph>`
- `EDGE.EXISTS <graph> <id>`

### `ANALYSIS` Commands

- `ANALYSIS.SHORTESTPATH <graph> <from_node> <to_node> [FORMAT simple|detailed]`
- `ANALYSIS.CENTRALITY <graph> <type> [node_id] [DIRECTION in|out|both]`
- `ANALYSIS.CLUSTERING <graph> [algorithm] [parameters_json]`
- `ANALYSIS.CYCLES <graph> [NODETYPE type1...] [EDGETYPE type1...] [FORMAT simple|detailed]`
- `ANALYSIS.TRAVERSE <graph> <start_node> [DIRECTION <dir>] [NODETYPES type1...] [EDGETYPES type1...] [FORMAT simple|detailed]`

*For detailed syntax, parameters, and examples for each command, please see the original `redis/README.md` file.*

## Storage Engine API (`storage.StorageEngine`)

The storage engine provides the core CRUD (Create, Read, Update, Delete) functionality for graphs, nodes, and edges.

### Graph Operations

- `CreateGraph(graph *models.Graph) error`
- `GetGraph(graphID models.GraphID) (*models.Graph, error)`
- `UpdateGraph(graph *models.Graph) error`
- `DeleteGraph(graphID models.GraphID) error`
- `ListGraphs() ([]*models.Graph, error)`

### Node Operations

- `CreateNode(graphID models.GraphID, node *models.Node) error`
- `GetNode(graphID models.GraphID, nodeID models.NodeID) (*models.Node, error)`
- `UpdateNode(graphID models.GraphID, node *models.Node) error`
- `DeleteNode(graphID models.GraphID, nodeID models.NodeID) error`
- `ListNodes(graphID models.GraphID) ([]*models.Node, error)`
- `ListNodesByType(graphID models.GraphID, nodeType models.NodeType) ([]*models.Node, error)`
- `FindNodesByAttribute(graphID models.GraphID, key string, value interface{}) ([]*models.Node, error)`

### Edge Operations

- `CreateEdge(graphID models.GraphID, edge *models.Edge) error`
- `GetEdge(graphID models.GraphID, edgeID models.EdgeID) (*models.Edge, error)`
- `UpdateEdge(graphID models.GraphID, edge *models.Edge) error`
- `DeleteEdge(graphID models.GraphID, edgeID models.EdgeID) error`
- `ListEdges(graphID models.GraphID) ([]*models.Edge, error)`
- `ListEdgesByType(graphID models.GraphID, edgeType models.EdgeType) ([]*models.Edge, error)`
- `GetOutgoingEdges(graphID models.GraphID, nodeID models.NodeID) ([]*models.Edge, error)`
- `GetIncomingEdges(graphID models.GraphID, nodeID models.NodeID) ([]*models.Edge, error)`
- `GetConnectedNodes(graphID models.GraphID, nodeID models.NodeID) ([]*models.Node, error)`

### Database Operations

- `Open(path string) error`
- `Close() error`
- `Backup(backupPath string) error`

## Analysis Engine API (`analysis.GraphAnalyzer`)

The analysis engine provides high-level functions for graph traversal, dependency analysis, and metrics calculation.

- `DepthFirstSearch(...)`
- `GetShortestPath(...)`
- `GetAllDependencies(...)`
- `GetAllDependents(...)`
- `HasCycles(...)`
- `GetGraphStats(...)`
- `GetRootNodes(...)`
- `GetLeafNodes(...)`
- `GetOrphanNodes(...)`
- `GetMaxDepth(...)`
- `GetConnectedComponentCount(...)`

## Docker (Production)

This project uses `docker-compose` to manage the multi-service production environment. The services are:
- `pathwaydb-redis`: The custom Redis-compatible server.
- `backend`: The Go WebSocket bridge.
- `frontend`: The React application served by Nginx.

### Prerequisites

- Docker
- Docker Compose

### Building and Running

To build and run the entire application stack, use the following command from the project root:

```bash
docker compose up --build
```

This will build the images for all services and start them in the correct order. The IDE will be available at `http://localhost:3000` by default.

### Configuration

You can configure the exposed ports by creating a `.env` file in the project root. If this file is not present, the default ports will be used.

**.env.example**
```
PORT=3000
REDIS_PORT=6379
WEBSOCKET_PORT=8081
```

### Stopping the Application

To stop and remove the containers, run:

```bash
docker compose down
```

## Testing

The project includes a comprehensive test suite covering the storage, analysis, and integration layers. It validates all CRUD operations, analysis algorithms, error handling, and edge cases.

To run the full test suite:

```bash
go test ./...
```

To run tests with coverage:

```bash
go test ./... -cover -v
```
