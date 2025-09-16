import React from 'react';
import {
  Box,
  Paper,
  Typography,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableRow,
  Chip
} from '@mui/material';
import { Circle, Timeline, Timer } from '@mui/icons-material';
import { GraphNode, GraphEdge } from '../types';

interface PropertiesPanelProps {
  selectedNode: GraphNode | null;
  selectedEdge: GraphEdge | null;
}

const timeUntil = (date: string) => {
  const now = new Date();
  const target = new Date(date);
  const diff = target.getTime() - now.getTime();

  if (diff <= 0) {
    return 'Expired';
  }

  const seconds = Math.floor(diff / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);

  if (days > 0) return `in ${days}d ${hours % 24}h`;
  if (hours > 0) return `in ${hours}h ${minutes % 60}m`;
  if (minutes > 0) return `in ${minutes}m ${seconds % 60}s`;
  return `in ${seconds}s`;
};

const PropertiesPanel: React.FC<PropertiesPanelProps> = ({
  selectedNode,
  selectedEdge
}) => {
  const renderNodeProperties = (node: GraphNode) => (
    <Box>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
        <Circle fontSize="small" sx={{ mr: 1, color: getNodeColor(node.type) }} />
        <Typography variant="h6" sx={{ fontSize: '14px', fontWeight: 'bold' }}>
          Node: {node.id}
        </Typography>
      </Box>
      
      <Chip 
        label={node.type} 
        size="small" 
        sx={{ 
          mb: 2,
          backgroundColor: getNodeColor(node.type),
          color: 'white'
        }}
      />

      <Typography variant="subtitle2" sx={{ mb: 1, fontSize: '12px', fontWeight: 'bold' }}>
        Properties
      </Typography>
      
      <TableContainer>
        <Table size="small">
          <TableBody>
            <TableRow>
              <TableCell sx={{ fontWeight: 'bold', fontSize: '11px', py: 0.5 }}>
                ID
              </TableCell>
              <TableCell sx={{ fontSize: '11px', py: 0.5 }}>
                {node.id}
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell sx={{ fontWeight: 'bold', fontSize: '11px', py: 0.5 }}>
                Type
              </TableCell>
              <TableCell sx={{ fontSize: '11px', py: 0.5 }}>
                {node.type}
              </TableCell>
            </TableRow>
            {node.expiresAt && (
              <TableRow>
                <TableCell sx={{ fontWeight: 'bold', fontSize: '11px', py: 0.5, display: 'flex', alignItems: 'center' }}>
                  <Timer fontSize="small" sx={{ mr: 0.5, color: 'text.secondary' }} /> TTL
                </TableCell>
                <TableCell sx={{ fontSize: '11px', py: 0.5 }}>
                  {timeUntil(node.expiresAt)}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </TableContainer>

      {Object.keys(node.attributes).length > 0 && (
        <>
          <Typography variant="subtitle2" sx={{ mb: 1, mt: 2, fontSize: '12px', fontWeight: 'bold' }}>
            Attributes
          </Typography>
          
          <TableContainer>
            <Table size="small">
              <TableBody>
                {Object.entries(node.attributes).map(([key, value]) => (
                  <TableRow key={key}>
                    <TableCell sx={{ fontWeight: 'bold', fontSize: '11px', py: 0.5 }}>
                      {key}
                    </TableCell>
                    <TableCell sx={{ fontSize: '11px', py: 0.5 }}>
                      {typeof value === 'object' ? JSON.stringify(value) : String(value)}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </>
      )}
    </Box>
  );

  const renderEdgeProperties = (edge: GraphEdge) => (
    <Box>
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
        <Timeline fontSize="small" sx={{ mr: 1, color: '#666' }} />
        <Typography variant="h6" sx={{ fontSize: '14px', fontWeight: 'bold' }}>
          Edge: {edge.id}
        </Typography>
      </Box>
      
      <Chip 
        label={edge.type} 
        size="small" 
        variant="outlined"
        sx={{ mb: 2 }}
      />

      <Typography variant="subtitle2" sx={{ mb: 1, fontSize: '12px', fontWeight: 'bold' }}>
        Properties
      </Typography>
      
      <TableContainer>
        <Table size="small">
          <TableBody>
            <TableRow>
              <TableCell sx={{ fontWeight: 'bold', fontSize: '11px', py: 0.5 }}>
                ID
              </TableCell>
              <TableCell sx={{ fontSize: '11px', py: 0.5 }}>
                {edge.id}
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell sx={{ fontWeight: 'bold', fontSize: '11px', py: 0.5 }}>
                Source
              </TableCell>
              <TableCell sx={{ fontSize: '11px', py: 0.5 }}>
                {edge.source}
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell sx={{ fontWeight: 'bold', fontSize: '11px', py: 0.5 }}>
                Target
              </TableCell>
              <TableCell sx={{ fontSize: '11px', py: 0.5 }}>
                {edge.target}
              </TableCell>
            </TableRow>
            <TableRow>
              <TableCell sx={{ fontWeight: 'bold', fontSize: '11px', py: 0.5 }}>
                Type
              </TableCell>
              <TableCell sx={{ fontSize: '11px', py: 0.5 }}>
                {edge.type}
              </TableCell>
            </TableRow>
            {edge.expiresAt && (
              <TableRow>
                <TableCell sx={{ fontWeight: 'bold', fontSize: '11px', py: 0.5, display: 'flex', alignItems: 'center' }}>
                  <Timer fontSize="small" sx={{ mr: 0.5, color: 'text.secondary' }} /> TTL
                </TableCell>
                <TableCell sx={{ fontSize: '11px', py: 0.5 }}>
                  {timeUntil(edge.expiresAt)}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </TableContainer>

      {Object.keys(edge.attributes).length > 0 && (
        <>
          <Typography variant="subtitle2" sx={{ mb: 1, mt: 2, fontSize: '12px', fontWeight: 'bold' }}>
            Attributes
          </Typography>
          
          <TableContainer>
            <Table size="small">
              <TableBody>
                {Object.entries(edge.attributes).map(([key, value]) => (
                  <TableRow key={key}>
                    <TableCell sx={{ fontWeight: 'bold', fontSize: '11px', py: 0.5 }}>
                      {key}
                    </TableCell>
                    <TableCell sx={{ fontSize: '11px', py: 0.5 }}>
                      {typeof value === 'object' ? JSON.stringify(value) : String(value)}
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </TableContainer>
        </>
      )}
    </Box>
  );

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

  return (
    <Paper sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <Box sx={{ p: 2, borderBottom: '1px solid #333' }}>
        <Typography variant="h6" sx={{ fontSize: '14px', fontWeight: 'bold' }}>
          Properties
        </Typography>
      </Box>

      <Box sx={{ flexGrow: 1, overflow: 'auto', p: 2 }}>
        {selectedNode && renderNodeProperties(selectedNode)}
        {selectedEdge && renderEdgeProperties(selectedEdge)}
        
        {!selectedNode && !selectedEdge && (
          <Box sx={{ textAlign: 'center', color: 'text.secondary', mt: 4 }}>
            <Typography variant="body2">
              Select a node or edge to view properties
            </Typography>
          </Box>
        )}
      </Box>
    </Paper>
  );
};

export default PropertiesPanel;
