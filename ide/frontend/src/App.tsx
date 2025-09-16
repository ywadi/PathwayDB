import React, { useState, useEffect, useCallback } from 'react';
import { Box, AppBar, Toolbar, Typography, Grid, Alert, Snackbar, IconButton, Tooltip } from '@mui/material';
import { Storage, MenuBook, ChevronLeft, ChevronRight } from '@mui/icons-material';
import GraphVisualization from './components/GraphVisualization';
import RedisConsole from './components/RedisConsole';
import Sidebar from './components/Sidebar';
import PropertiesPanel from './components/PropertiesPanel';
import DocumentationPage from './components/DocumentationPage';
import { RedisWebSocket } from './services/RedisWebSocket';
import { Graph, GraphNode, GraphEdge, ConnectionStatus, RedisResponse } from './types';

const App: React.FC = () => {
  const [redisClient] = useState(() => new RedisWebSocket('localhost', 8081));
  const [connectionStatus, setConnectionStatus] = useState<ConnectionStatus>({
    connected: false,
    host: 'localhost',
    port: 8081
  });
  const [graphs, setGraphs] = useState<Graph[]>([]);
  const [selectedGraph, setSelectedGraph] = useState<string | null>(null);
  const [selectedNode, setSelectedNode] = useState<GraphNode | null>(null);
  const [selectedEdge, setSelectedEdge] = useState<GraphEdge | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [showDocs, setShowDocs] = useState<boolean>(false);
  const [propertiesHeight, setPropertiesHeight] = useState<number>(30); // Percentage
  const [sidebarCollapsed, setSidebarCollapsed] = useState<boolean>(false);

  // Initialize WebSocket connection
  useEffect(() => {
    redisClient.onConnectionChange = setConnectionStatus;
    redisClient.onError = setError;

    const connect = async () => {
      try {
        await redisClient.connect();
        await loadGraphs();
      } catch (err) {
        setError('Failed to connect to PathwayDB Redis server');
      }
    };

    connect();

    return () => {
      redisClient.disconnect();
    };
  }, []);

  const loadGraphs = useCallback(async () => {
    if (!redisClient.isConnected()) return;
    try {
      const graphList = await redisClient.listGraphs();
      const loadedGraphs: Graph[] = [];

      // Process graph list (comes as [id, description, id, description, ...])
      for (let i = 0; i < graphList.length; i += 2) {
        const graphId = graphList[i];
        const description = graphList[i + 1] || '';

        try {
          // Load nodes and edges for each graph
          const [nodeList, edgeList] = await Promise.all([
            redisClient.listNodes(graphId).catch(() => []),
            redisClient.listEdges(graphId).catch(() => [])
          ]);

          // Convert to proper format
          const graphNodes: GraphNode[] = [];
          const graphEdges: GraphEdge[] = [];

          // Process nodes - NODE.LIST returns [id, type, id, type, ...]
          for (let j = 0; j < nodeList.length; j += 2) {
            if (nodeList[j]) {
              const nodeId = nodeList[j];
              const nodeType = nodeList[j + 1] || 'default';
              
              // Get detailed node info
              try {
                const nodeDetails = await redisClient.getNode(graphId, nodeId);
                let attributes = {};
                let expiresAt: string | undefined = undefined;

                if (nodeDetails && nodeDetails.length >= 3) {
                  try {
                    const attrStr = nodeDetails[2] || '{}';
                    if (attrStr.startsWith('{') && attrStr.endsWith('}')) {
                      attributes = JSON.parse(attrStr);
                    }
                  } catch (e) {
                    console.warn(`Failed to parse node attributes for ${nodeId}:`, nodeDetails[2]);
                  }
                }

                if (nodeDetails && nodeDetails.length >= 4 && nodeDetails[3]) {
                  expiresAt = nodeDetails[3];
                }
                
                graphNodes.push({
                  id: nodeId,
                  type: nodeType,
                  attributes,
                  expiresAt
                });
              } catch (e) {
                // If we can't get details, add basic node
                graphNodes.push({
                  id: nodeId,
                  type: nodeType,
                  attributes: {}
                });
              }
            }
          }

          // Process edges - EDGE.LIST returns [id, source, target, type, ...]
          for (let j = 0; j < edgeList.length; j += 4) {
            if (edgeList[j]) {
              const edgeId = edgeList[j];
              const source = edgeList[j + 1];
              const target = edgeList[j + 2];
              const edgeType = edgeList[j + 3] || 'default';
              
              // Get detailed edge info
              try {
                const edgeDetails = await redisClient.getEdge(graphId, edgeId);
                let attributes = {};
                let expiresAt: string | undefined = undefined;

                if (edgeDetails && edgeDetails.length >= 5) {
                  try {
                    const attrStr = edgeDetails[4] || '{}';
                    if (attrStr.startsWith('{') && attrStr.endsWith('}')) {
                      attributes = JSON.parse(attrStr);
                    }
                  } catch (e) {
                    console.warn(`Failed to parse edge attributes for ${edgeId}:`, edgeDetails[4]);
                  }
                }

                if (edgeDetails && edgeDetails.length >= 6 && edgeDetails[5]) {
                  expiresAt = edgeDetails[5];
                }
                
                graphEdges.push({
                  id: edgeId,
                  source,
                  target,
                  type: edgeType,
                  attributes,
                  expiresAt
                });
              } catch (e) {
                // If we can't get details, add basic edge
                graphEdges.push({
                  id: edgeId,
                  source,
                  target,
                  type: edgeType,
                  attributes: {}
                });
              }
            }
          }

          // Check for cycles
          let hasCycles = false;
          try {
            const cycleResponse = await redisClient.executeCommand('ANALYSIS.CYCLES', [graphId]);
            hasCycles = cycleResponse.value && Array.isArray(cycleResponse.value) && cycleResponse.value.length > 0;
          } catch (err) {
            console.warn(`Failed to check cycles for graph ${graphId}:`, err);
          }

          // Add graph with loaded data
          loadedGraphs.push({
            id: graphId,
            name: graphId,
            description,
            nodes: graphNodes,
            edges: graphEdges,
            hasCycles
          });
        } catch (err) {
          console.warn(`Failed to load data for graph ${graphId}:`, err);
          // Add graph with empty data
          loadedGraphs.push({
            id: graphId,
            name: graphId,
            description,
            nodes: [],
            edges: []
          });
        }
      }

      setGraphs(loadedGraphs);
      
      // Auto-select first graph if none selected
      if (!selectedGraph && loadedGraphs.length > 0) {
        setSelectedGraph(loadedGraphs[0].id);
      }
    } catch (err) {
      setError('Failed to load graphs');
      console.error('Error loading graphs:', err);
    }
  }, [redisClient, selectedGraph]);

  const handleExecuteCommand = async (command: string): Promise<RedisResponse> => {
    const parts = command.trim().split(/\s+/);
    const cmd = parts[0].toUpperCase();
    const args = parts.slice(1);

    const response = await redisClient.executeCommand(cmd, args);
    
    // Refresh graphs if command might have changed data
    if (cmd.startsWith('GRAPH.') || cmd.startsWith('NODE.') || cmd.startsWith('EDGE.')) {
      setTimeout(loadGraphs, 100);
    }

    return response;
  };

  const handleGraphSelect = (graphId: string) => {
    setSelectedGraph(graphId);
    setSelectedNode(null);
    setSelectedEdge(null);
  };

  const handleNodeSelect = useCallback((node: GraphNode | null) => {
    setSelectedNode(node);
    setSelectedEdge(null);
  }, []);

  const handleEdgeSelect = useCallback((edge: GraphEdge | null) => {
    setSelectedEdge(edge);
    setSelectedNode(null);
  }, []);

  const handleOpenDocumentation = () => {
    setShowDocs(true);
  };

  const handleCloseDocumentation = useCallback(() => {
    // Clear selections to prevent Cytoscape.js from trying to access stale references
    setSelectedNode(null);
    setSelectedEdge(null);
    
    // Force a complete re-render by clearing the graph selection temporarily
    const currentGraphId = selectedGraph;
    setSelectedGraph(null);
    
    // Use a longer timeout to ensure complete cleanup
    setTimeout(() => {
      setShowDocs(false);
      // Restore the graph selection after the view has switched
      setTimeout(() => {
        if (currentGraphId) {
          setSelectedGraph(currentGraphId);
        }
      }, 100);
    }, 50);
  }, [selectedGraph]);

  const currentGraph = graphs.find(g => g.id === selectedGraph);

  if (showDocs) {
    return <DocumentationPage onClose={handleCloseDocumentation} />;
  }

  return (
    <Box sx={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      {/* Header */}
      <AppBar position="static" sx={{ zIndex: 1201 }}>
        <Toolbar variant="dense" sx={{ minHeight: '48px' }}>
          <Storage sx={{ mr: 1 }} />
          <Typography variant="h6" sx={{ flexGrow: 1, fontSize: '16px' }}>
            PathwayDB IDE
          </Typography>
          <Tooltip title="Open Documentation">
            <IconButton
              color="inherit"
              onClick={handleOpenDocumentation}
              sx={{ mr: 2 }}
            >
              <MenuBook />
            </IconButton>
          </Tooltip>
          <Typography variant="body2" sx={{ fontSize: '12px' }}>
            {connectionStatus.connected ? 
              `Connected to ${connectionStatus.host}:${connectionStatus.port}` : 
              'Disconnected'
            }
          </Typography>
        </Toolbar>
      </AppBar>

      {/* Main Content */}
      <Box sx={{ flexGrow: 1, display: 'flex', overflow: 'hidden' }}>
        <Grid container sx={{ height: '100%' }}>
          {/* Sidebar */}
          {!sidebarCollapsed && (
            <Grid item xs={2.5} sx={{ borderRight: '1px solid #333' }}>
              <Sidebar
                graphs={graphs}
                selectedGraph={selectedGraph}
                onGraphSelect={handleGraphSelect}
                onRefresh={loadGraphs}
                connected={connectionStatus.connected}
              />
            </Grid>
          )}

          {/* Sidebar Toggle Button */}
          <Box
            sx={{
              position: 'absolute',
              left: sidebarCollapsed ? 0 : 'calc(20.833% - 16px)', // 2.5/12 = 20.833%
              top: '50%',
              transform: 'translateY(-50%)',
              zIndex: 1000,
              backgroundColor: '#333',
              borderRadius: '0 8px 8px 0',
              borderLeft: sidebarCollapsed ? 'none' : '1px solid #555',
              borderTop: '1px solid #555',
              borderRight: '1px solid #555',
              borderBottom: '1px solid #555',
            }}
          >
            <IconButton
              size="small"
              onClick={() => setSidebarCollapsed(!sidebarCollapsed)}
              sx={{
                color: '#fff',
                p: 0.5,
                '&:hover': {
                  backgroundColor: '#58a6ff',
                }
              }}
            >
              {sidebarCollapsed ? <ChevronRight /> : <ChevronLeft />}
            </IconButton>
          </Box>

          {/* Main Visualization Area */}
          <Grid item xs={sidebarCollapsed ? 8.5 : 6} sx={{ display: 'flex', flexDirection: 'column' }}>
            <Box sx={{ flexGrow: 1, borderRight: '1px solid #333' }}>
              <GraphVisualization
                nodes={currentGraph?.nodes || []}
                edges={currentGraph?.edges || []}
                onNodeSelect={handleNodeSelect}
                onEdgeSelect={handleEdgeSelect}
                onRefresh={loadGraphs}
              />
            </Box>
          </Grid>

          {/* Right Panel */}
          <Grid item xs={3.5} sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
            {/* Properties Panel */}
            <Box sx={{ 
              height: `${propertiesHeight}%`, 
              minHeight: '150px',
              borderBottom: '1px solid #333',
              overflow: 'hidden'
            }}>
              <PropertiesPanel
                selectedNode={selectedNode}
                selectedEdge={selectedEdge}
              />
            </Box>

            {/* Resize Handle */}
            <Box
              sx={{
                height: '4px',
                backgroundColor: '#333',
                cursor: 'row-resize',
                borderTop: '1px solid #555',
                borderBottom: '1px solid #555',
                '&:hover': {
                  backgroundColor: '#58a6ff',
                }
              }}
              onMouseDown={(e) => {
                const startY = e.clientY;
                const startHeight = propertiesHeight;
                const containerHeight = e.currentTarget.parentElement?.clientHeight || 0;

                const handleMouseMove = (e: MouseEvent) => {
                  const deltaY = e.clientY - startY;
                  const deltaPercent = (deltaY / containerHeight) * 100;
                  const newHeight = Math.min(Math.max(startHeight + deltaPercent, 20), 70);
                  setPropertiesHeight(newHeight);
                };

                const handleMouseUp = () => {
                  document.removeEventListener('mousemove', handleMouseMove);
                  document.removeEventListener('mouseup', handleMouseUp);
                };

                document.addEventListener('mousemove', handleMouseMove);
                document.addEventListener('mouseup', handleMouseUp);
              }}
            />

            {/* Console */}
            <Box sx={{ 
              height: `${100 - propertiesHeight - 1}%`, 
              minHeight: '300px',
              display: 'flex',
              flexDirection: 'column',
              overflow: 'hidden'
            }}>
              <RedisConsole
                onExecuteCommand={handleExecuteCommand}
                connected={connectionStatus.connected}
              />
            </Box>
          </Grid>
        </Grid>
      </Box>

      {/* Error Snackbar */}
      <Snackbar
        open={!!error}
        autoHideDuration={6000}
        onClose={() => setError(null)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
      >
        <Alert onClose={() => setError(null)} severity="error" sx={{ width: '100%' }}>
          {error}
        </Alert>
      </Snackbar>
    </Box>
  );
};

export default App;
