import React, { useState } from 'react';
import {
  Box,
  Paper,
  Typography,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  ListItemButton,
  Collapse,
  IconButton,
  Tooltip,
  Chip,
  Divider,
  TextField,
  InputAdornment
} from '@mui/material';
import {
  ExpandLess,
  ExpandMore,
  Storage,
  Circle,
  Timeline,
  Refresh,
  Loop,
  Search,
  Clear
} from '@mui/icons-material';
import { Graph, GraphNode, GraphEdge } from '../types';

interface SidebarProps {
  graphs: Graph[];
  selectedGraph: string | null;
  onGraphSelect: (graphId: string) => void;
  onRefresh: () => void;
  connected: boolean;
}

const Sidebar: React.FC<SidebarProps> = ({
  graphs,
  selectedGraph,
  onGraphSelect,
  onRefresh,
  connected
}) => {
  const [expandedGraphs, setExpandedGraphs] = useState<Set<string>>(new Set());
  const [searchQuery, setSearchQuery] = useState<string>('');

  const handleToggleExpand = (graphId: string) => {
    setExpandedGraphs(prev => {
      const newSet = new Set(prev);
      if (newSet.has(graphId)) {
        newSet.delete(graphId);
      } else {
        newSet.add(graphId);
      }
      return newSet;
    });
  };

  const getNodeTypeCount = (nodes: GraphNode[]) => {
    const counts: Record<string, number> = {};
    nodes.forEach(node => {
      counts[node.type] = (counts[node.type] || 0) + 1;
    });
    return counts;
  };

  const getEdgeTypeCount = (edges: GraphEdge[]) => {
    const counts: Record<string, number> = {};
    edges.forEach(edge => {
      counts[edge.type] = (counts[edge.type] || 0) + 1;
    });
    return counts;
  };

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

  // Filter graphs based on search query
  const filteredGraphs = graphs.filter(graph => 
    graph.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    (graph.description && graph.description.toLowerCase().includes(searchQuery.toLowerCase()))
  );

  return (
    <Paper sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <Box sx={{ p: 1.5, borderBottom: '1px solid #333' }}>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <Typography variant="h6" sx={{ fontSize: '13px', fontWeight: 600 }}>
            Graph Explorer
          </Typography>
          <Tooltip title="Refresh">
            <span>
              <IconButton size="small" onClick={onRefresh} disabled={!connected} sx={{ p: 0.5 }}>
                <Refresh sx={{ fontSize: 16 }} />
              </IconButton>
            </span>
          </Tooltip>
        </Box>
        
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mt: 0.5 }}>
          <Chip 
            label={`${filteredGraphs.length}/${graphs.length} graphs`} 
            size="small" 
            color="primary" 
            variant="outlined"
            sx={{ height: 20, fontSize: '10px' }}
          />
        </Box>
        
        {/* Search Field */}
        <TextField
          fullWidth
          size="small"
          placeholder="Search graphs..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          sx={{
            mt: 1,
            '& .MuiOutlinedInput-root': {
              height: '28px',
              fontSize: '11px',
              backgroundColor: '#2a2a2a',
              '& fieldset': {
                borderColor: '#444',
              },
              '&:hover fieldset': {
                borderColor: '#666',
              },
              '&.Mui-focused fieldset': {
                borderColor: '#58a6ff',
              },
            },
            '& .MuiInputBase-input': {
              padding: '4px 8px',
              color: '#e6edf3',
              '&::placeholder': {
                color: '#8b949e',
                opacity: 1,
              },
            },
          }}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <Search sx={{ fontSize: 14, color: '#8b949e' }} />
              </InputAdornment>
            ),
            endAdornment: searchQuery && (
              <InputAdornment position="end">
                <IconButton
                  size="small"
                  onClick={() => setSearchQuery('')}
                  sx={{ p: 0.25, color: '#8b949e' }}
                >
                  <Clear sx={{ fontSize: 12 }} />
                </IconButton>
              </InputAdornment>
            ),
          }}
        />
      </Box>

      <Box sx={{ flexGrow: 1, overflow: 'auto' }}>
        {graphs.length === 0 ? (
          <Box sx={{ p: 2, textAlign: 'center' }}>
            <Typography variant="body2" color="text.secondary">
              {connected ? 'No graphs found' : 'Not connected'}
            </Typography>
            {connected && (
              <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mt: 1 }}>
                Create a graph using: GRAPH.CREATE &lt;name&gt;
              </Typography>
            )}
          </Box>
        ) : filteredGraphs.length === 0 ? (
          <Box sx={{ p: 2, textAlign: 'center' }}>
            <Typography variant="body2" color="text.secondary">
              No graphs match "{searchQuery}"
            </Typography>
            <Typography variant="caption" color="text.secondary" sx={{ display: 'block', mt: 1 }}>
              Try a different search term
            </Typography>
          </Box>
        ) : (
          <List dense sx={{ py: 0 }}>
            {filteredGraphs.map((graph) => (
              <React.Fragment key={graph.id}>
                <ListItem disablePadding>
                  <ListItemButton
                    selected={selectedGraph === graph.id}
                    onClick={() => onGraphSelect(graph.id)}
                    sx={{ py: 0.5, px: 1 }}
                  >
                    <ListItemIcon sx={{ minWidth: 28 }}>
                      <Storage sx={{ fontSize: 16 }} />
                    </ListItemIcon>
                    <ListItemText
                      primary={
                        <Typography variant="body2" sx={{ fontWeight: 500, fontSize: '12px' }}>
                          {graph.name}
                        </Typography>
                      }
                      secondary={
                        <Typography variant="caption" color="text.secondary" sx={{ fontSize: '10px' }}>
                          {graph.nodes.length}n, {graph.edges.length}e
                        </Typography>
                      }
                    />
                    <IconButton
                      size="small"
                      onClick={(e) => {
                        e.stopPropagation();
                        handleToggleExpand(graph.id);
                      }}
                      sx={{ p: 0.25 }}
                    >
                      {expandedGraphs.has(graph.id) ? (
                        <ExpandLess sx={{ fontSize: 16 }} />
                      ) : (
                        <ExpandMore sx={{ fontSize: 16 }} />
                      )}
                    </IconButton>
                  </ListItemButton>
                </ListItem>

                <Collapse in={expandedGraphs.has(graph.id)} timeout="auto" unmountOnExit>
                  <Box sx={{ pl: 1.5 }}>
                    {/* Node Types */}
                    <ListItem sx={{ py: 0.25 }}>
                      <ListItemIcon sx={{ minWidth: 24 }}>
                        <Circle sx={{ fontSize: 12 }} />
                      </ListItemIcon>
                      <ListItemText
                        primary={
                          <Typography variant="body2" sx={{ fontSize: '11px', fontWeight: 500 }}>
                            Nodes ({graph.nodes.length})
                          </Typography>
                        }
                      />
                    </ListItem>
                    
                    {Object.entries(getNodeTypeCount(graph.nodes)).map(([type, count]) => (
                      <ListItem key={type} sx={{ pl: 4, py: 0.125 }}>
                        <Box
                          sx={{
                            width: 6,
                            height: 6,
                            borderRadius: '50%',
                            backgroundColor: getNodeColor(type),
                            mr: 0.75,
                            flexShrink: 0
                          }}
                        />
                        <Typography variant="caption" sx={{ flexGrow: 1, fontSize: '10px' }}>
                          {type}
                        </Typography>
                        <Typography variant="caption" color="text.secondary" sx={{ fontSize: '10px' }}>
                          {count}
                        </Typography>
                      </ListItem>
                    ))}

                    {/* Edge Types */}
                    <ListItem sx={{ py: 0.25 }}>
                      <ListItemIcon sx={{ minWidth: 24 }}>
                        <Timeline sx={{ fontSize: 12 }} />
                      </ListItemIcon>
                      <ListItemText
                        primary={
                          <Typography variant="body2" sx={{ fontSize: '11px', fontWeight: 500 }}>
                            Edges ({graph.edges.length})
                          </Typography>
                        }
                      />
                    </ListItem>
                    
                    {Object.entries(getEdgeTypeCount(graph.edges)).map(([type, count]) => (
                      <ListItem key={type} sx={{ pl: 4, py: 0.125 }}>
                        <Box
                          sx={{
                            width: 6,
                            height: 1.5,
                            backgroundColor: '#666',
                            mr: 0.75,
                            flexShrink: 0
                          }}
                        />
                        <Typography variant="caption" sx={{ flexGrow: 1, fontSize: '10px' }}>
                          {type}
                        </Typography>
                        <Typography variant="caption" color="text.secondary" sx={{ fontSize: '10px' }}>
                          {count}
                        </Typography>
                      </ListItem>
                    ))}

                    {/* Cycles Detection */}
                    <ListItem sx={{ py: 0.25 }}>
                      <ListItemIcon sx={{ minWidth: 24 }}>
                        <Loop sx={{ fontSize: 12 }} />
                      </ListItemIcon>
                      <ListItemText
                        primary={
                          <Typography variant="body2" sx={{ fontSize: '11px', fontWeight: 500 }}>
                            Status
                          </Typography>
                        }
                      />
                    </ListItem>
                    
                    <ListItem sx={{ pl: 4, py: 0.125 }}>
                      <Box
                        sx={{
                          width: 6,
                          height: 6,
                          borderRadius: '50%',
                          backgroundColor: graph.hasCycles ? '#f44336' : '#4caf50',
                          mr: 0.75,
                          flexShrink: 0
                        }}
                      />
                      <Typography variant="caption" sx={{ flexGrow: 1, fontSize: '10px' }}>
                        {graph.hasCycles ? 'Cyclic' : 'Acyclic'}
                      </Typography>
                    </ListItem>
                  </Box>
                </Collapse>

                <Divider sx={{ mx: 1 }} />
              </React.Fragment>
            ))}
          </List>
        )}
      </Box>
    </Paper>
  );
};

export default Sidebar;
