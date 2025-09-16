# PathwayDB IDE

A modern, professional web-based IDE for PathwayDB Graph Database with real-time graph visualization and Redis command console.

## Features

- **Professional UI**: Modern dark theme with compact, business-ready design
- **Graph Visualization**: Interactive graph visualization using Cytoscape.js with multiple layouts
- **Redis Console**: Real-time command execution with syntax highlighting and command history
- **WebSocket Integration**: Direct connection to Redis protocol for real-time updates
- **Graph Explorer**: Browse graphs, nodes, and edges with detailed statistics
- **Properties Panel**: Inspect selected nodes and edges with full attribute details and TTL expiration times.
- **Documentation Viewer**: In-app viewer for all project and command documentation.

## Architecture

- **Frontend**: React + TypeScript + Material-UI + Cytoscape.js
- **Backend**: Go WebSocket server with Redis protocol bridge
- **Communication**: WebSocket for real-time Redis command execution
- **Visualization**: Cytoscape.js with force-directed and hierarchical layouts

## Quick Start

The easiest way to get started is with the integrated start script. This will handle all dependencies, build the necessary components, and launch the Redis server, backend, and frontend.

```bash
# From the project root
./ide/start.sh
```

Once running, you can access the IDE at **http://localhost:3000**.

For manual development, see the `Development` section below.

## Usage

For a complete reference of all supported Redis commands, including syntax and examples, please see the [Command Reference](./COMMANDS.md).

### UI Components

#### Graph Visualization
- **Pan/Zoom**: Mouse controls for navigation
- **Layouts**: Toggle between force-directed and hierarchical layouts
- **Node Selection**: Click nodes/edges to view properties
- **Auto-fit**: Automatically fit graph to viewport

#### Graph Explorer
- **Graph List**: Browse all available graphs
- **Statistics**: View node/edge counts and types
- **Expandable Details**: See node and edge type breakdowns

#### Properties Panel
- **Node Properties**: View selected node attributes and metadata
- **Edge Properties**: View selected edge relationships and attributes
- **Real-time Updates**: Properties update as you select different elements

#### Redis Console
- **Command History**: Use ↑/↓ arrows to navigate command history
- **Syntax Highlighting**: Commands are highlighted for better readability
- **Real-time Execution**: Commands execute immediately with live feedback
- **Connection Status**: Visual indicator of Redis connection status

## Configuration

### Backend Configuration

The WebSocket bridge server accepts these flags:

```bash
go run main.go -addr :8080 -redis localhost:6379
```

- `-addr`: WebSocket server address (default: :8080)
- `-redis`: Redis server address (default: localhost:6379)

### Frontend Configuration

The frontend connects to the WebSocket server at `ws://localhost:8080/ws` by default. To change this, modify the `RedisWebSocket` constructor in `src/services/RedisWebSocket.ts`.

## Development

### Project Structure

```
ide/
├── frontend/                 # React frontend application
│   ├── src/
│   │   ├── components/      # React components
│   │   ├── services/        # WebSocket service
│   │   ├── types/          # TypeScript type definitions
│   │   └── App.tsx         # Main application component
│   ├── public/             # Static assets
│   └── package.json        # Frontend dependencies
├── backend/                 # Go WebSocket bridge server
│   ├── main.go             # WebSocket server implementation
│   └── go.mod              # Backend dependencies
└── README.md               # This file
```

### Adding New Features

1. **New Redis Commands**: Add command handlers in the backend Redis proxy
2. **UI Components**: Create new React components in `src/components/`
3. **Graph Layouts**: Add new Cytoscape.js layouts in `GraphVisualization.tsx`
4. **Styling**: Modify Material-UI theme in `src/index.tsx`

## Troubleshooting

### Connection Issues

1. **Redis Server Not Running**: Ensure PathwayDB Redis server is running on port 6379
2. **WebSocket Connection Failed**: Check that the IDE backend is running on port 8080
3. **CORS Errors**: The backend allows all origins for development

### Performance Issues

1. **Large Graphs**: Use graph filtering or pagination for graphs with many nodes
2. **Memory Usage**: Close unused browser tabs and restart the IDE if needed
3. **Network Latency**: Ensure Redis server and IDE backend are on the same network

### Development Issues

1. **Hot Reload**: Frontend supports hot reload; backend requires restart for changes
2. **TypeScript Errors**: Run `npm run build` to check for type errors
3. **Go Dependencies**: Run `go mod tidy` if you encounter import issues

## License

This IDE is part of the PathwayDB project and follows the same licensing terms.
