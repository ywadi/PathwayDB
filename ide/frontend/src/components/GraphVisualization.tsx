import React, { useEffect, useRef, useState, useCallback } from 'react';
import { Box, Paper, IconButton, Toolbar, Typography, Tooltip } from '@mui/material';
import { 
  ZoomIn, 
  ZoomOut, 
  CenterFocusStrong, 
  GridOn,
  AccountTree,
  Refresh
} from '@mui/icons-material';
import cytoscape, { Core, ElementDefinition } from 'cytoscape';
import dagre from 'cytoscape-dagre';
import fcose from 'cytoscape-fcose';
import { GraphNode, GraphEdge } from '../types';

// Register layout extensions
cytoscape.use(dagre);
cytoscape.use(fcose);

interface GraphVisualizationProps {
  nodes: GraphNode[];
  edges: GraphEdge[];
  onNodeSelect?: (node: GraphNode | null) => void;
  onEdgeSelect?: (edge: GraphEdge | null) => void;
  onRefresh?: () => void;
}

const GraphVisualization: React.FC<GraphVisualizationProps> = ({
  nodes,
  edges,
  onNodeSelect,
  onEdgeSelect,
  onRefresh
}) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const cyRef = useRef<Core | null>(null);
  const [layout, setLayout] = useState<'force' | 'hierarchical'>('force');
  
  // Cleanup function to properly destroy Cytoscape instance
  const cleanupCytoscape = useCallback(() => {
    if (cyRef.current && !cyRef.current.destroyed()) {
      try {
        // Stop all running animations and layouts
        cyRef.current.stop();
        // Remove all event listeners
        cyRef.current.removeAllListeners();
        // Clear all elements
        cyRef.current.elements().remove();
        // Destroy the instance
        cyRef.current.destroy();
      } catch (error) {
        console.warn('Error during Cytoscape cleanup:', error);
      }
      cyRef.current = null;
    }
  }, []);

  // Add cleanup on component unmount
  useEffect(() => {
    return cleanupCytoscape;
  }, [cleanupCytoscape]);

  useEffect(() => {
    if (!containerRef.current) return;

    // Initialize Cytoscape
    const cy = cytoscape({
      container: containerRef.current,
      style: [
        {
          selector: 'node',
          style: {
            'background-color': (ele: any) => getNodeColor(ele.data('type')),
            'label': 'data(id)',
            'text-valign': 'center',
            'text-halign': 'center',
            'font-size': '10px',
            'font-weight': 'bold',
            'color': 'white',
            'text-outline-width': 1,
            'text-outline-color': '#000',
            'width': 60,
            'height': 60,
            'border-width': 2,
            'border-color': '#fff',
            'border-opacity': 0.8
          }
        },
        {
          selector: 'edge',
          style: {
            'width': 2,
            'line-color': '#666',
            'target-arrow-color': '#666',
            'target-arrow-shape': 'triangle',
            'curve-style': 'bezier',
            'label': 'data(id)',
            'font-size': '8px',
            'text-rotation': 'autorotate',
            'text-margin-y': -10,
            'color': '#999'
          }
        },
        {
          selector: ':selected',
          style: {
            'border-width': 3,
            'border-color': '#00bcd4',
            'line-color': '#00bcd4',
            'target-arrow-color': '#00bcd4',
            'source-arrow-color': '#00bcd4'
          }
        }
      ],
      layout: {
        name: 'preset',
        padding: 50
      } as any
    });

    cyRef.current = cy;

    return () => {
      // Use the cleanup function
      cleanupCytoscape();
    };
  }, []);

  // Separate useEffect for event handlers - only update when handlers change
  useEffect(() => {
    if (!cyRef.current) return;

    const cy = cyRef.current;
    
    // Remove existing handlers
    cy.off('tap');

    // Event handlers
    cy.on('tap', 'node', (evt) => {
      const node = evt.target;
      const nodeData = nodes.find(n => n.id === node.id());
      console.log('Node clicked:', node.id(), 'Found data:', nodeData);
      onNodeSelect?.(nodeData || null);
    });

    cy.on('tap', 'edge', (evt) => {
      const edge = evt.target;
      const edgeData = edges.find(e => e.id === edge.id());
      console.log('Edge clicked:', edge.id(), 'Found data:', edgeData);
      onEdgeSelect?.(edgeData || null);
    });

    cy.on('tap', (evt) => {
      if (evt.target === cy) {
        cy.$(':selected').unselect();
        onNodeSelect?.(null);
        onEdgeSelect?.(null);
      }
    });
  }, [nodes, edges, onNodeSelect, onEdgeSelect]);

  useEffect(() => {
    if (!cyRef.current) return;

    // Convert data to Cytoscape format
    const elements: ElementDefinition[] = [
      ...nodes.map(node => ({
        data: {
          id: node.id,
          type: node.type,
          ...node.attributes
        },
        position: node.position
      })),
      ...edges.map(edge => ({
        data: {
          id: edge.id,
          source: edge.source,
          target: edge.target,
          type: edge.type,
          ...edge.attributes
        }
      }))
    ];

    // Update graph
    cyRef.current.elements().remove();
    cyRef.current.add(elements);
    
    // Apply layout
    const layoutName = layout === 'hierarchical' ? 'dagre' : 'cose';
    const layoutOptions: any = {
      name: layoutName,
      animate: true,
      animationDuration: 1000,
      fit: true,
      padding: 50,
      ...(layoutName === 'dagre' && {
        rankDir: 'TB',
        nodeSep: 50,
        edgeSep: 10,
        rankSep: 50
      })
    };
    
    cyRef.current.layout(layoutOptions).run();

  }, [nodes, edges, layout]);

  const getNodeColor = (type: string): string => {
    const colors: Record<string, string> = {
      'service': '#4caf50',
      'database': '#f44336',
      'cache': '#ff9800',
      'api': '#2196f3',
      'user': '#9c27b0',
      'default': '#757575'
    };
    return colors[type] || colors.default;
  };

  const handleZoomIn = () => {
    cyRef.current?.zoom(cyRef.current.zoom() * 1.2);
  };

  const handleZoomOut = () => {
    cyRef.current?.zoom(cyRef.current.zoom() * 0.8);
  };

  const handleFit = () => {
    cyRef.current?.fit(undefined, 50);
  };

  const handleToggleLayout = () => {
    setLayout(prev => prev === 'force' ? 'hierarchical' : 'force');
  };

  return (
    <Paper sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <Toolbar variant="dense" sx={{ minHeight: '48px', borderBottom: '1px solid #333' }}>
        <Typography variant="h6" sx={{ flexGrow: 1, fontSize: '14px' }}>
          Graph Visualization ({nodes.length} nodes, {edges.length} edges)
        </Typography>
        
        <Tooltip title="Refresh">
          <IconButton size="small" onClick={onRefresh}>
            <Refresh fontSize="small" />
          </IconButton>
        </Tooltip>
        
        <Tooltip title="Toggle Layout">
          <IconButton size="small" onClick={handleToggleLayout}>
            {layout === 'force' ? <AccountTree fontSize="small" /> : <GridOn fontSize="small" />}
          </IconButton>
        </Tooltip>
        
        <Tooltip title="Zoom In">
          <IconButton size="small" onClick={handleZoomIn}>
            <ZoomIn fontSize="small" />
          </IconButton>
        </Tooltip>
        
        <Tooltip title="Zoom Out">
          <IconButton size="small" onClick={handleZoomOut}>
            <ZoomOut fontSize="small" />
          </IconButton>
        </Tooltip>
        
        <Tooltip title="Fit to Screen">
          <IconButton size="small" onClick={handleFit}>
            <CenterFocusStrong fontSize="small" />
          </IconButton>
        </Tooltip>
      </Toolbar>
      
      <Box 
        ref={containerRef}
        sx={{ 
          flexGrow: 1, 
          backgroundColor: '#0a0a0a',
          position: 'relative'
        }}
      />
    </Paper>
  );
};

export default GraphVisualization;
