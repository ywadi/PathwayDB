import React, { useState, useRef, useEffect } from 'react';
import { 
  Box, 
  Paper, 
  Typography, 
  TextField, 
  Toolbar, 
  IconButton,
  List,
  ListItem,
  ListItemText,
  Chip
} from '@mui/material';
import { 
  Send, 
  Clear, 
  Terminal
} from '@mui/icons-material';
import { ConsoleEntry, RedisResponse } from '../types';

interface RedisConsoleProps {
  onExecuteCommand: (command: string) => Promise<RedisResponse>;
  connected: boolean;
}

const RedisConsole: React.FC<RedisConsoleProps> = ({ onExecuteCommand, connected }) => {
  const [command, setCommand] = useState('');
  const [history, setHistory] = useState<ConsoleEntry[]>([]);
  const [historyIndex, setHistoryIndex] = useState(-1);
  const [commandHistory, setCommandHistory] = useState<string[]>([]);
  const listRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // Scroll to bottom when new entries are added
    if (listRef.current) {
      listRef.current.scrollTop = listRef.current.scrollHeight;
    }
  }, [history]);

  const handleExecuteCommand = async () => {
    if (!command.trim() || !connected) return;

    const entry: ConsoleEntry = {
      id: Date.now().toString(),
      command: command.trim(),
      timestamp: Date.now(),
      status: 'pending'
    };

    setHistory(prev => [...prev, entry]);
    setCommandHistory(prev => [command.trim(), ...prev.slice(0, 49)]); // Keep last 50 commands
    setCommand('');
    setHistoryIndex(-1);

    try {
      const response = await onExecuteCommand(command.trim());
      
      setHistory(prev => prev.map(h => 
        h.id === entry.id 
          ? { ...h, response, status: 'success' }
          : h
      ));
    } catch (error) {
      const errorResponse: RedisResponse = {
        type: 'error',
        value: error instanceof Error ? error.message : 'Unknown error',
        timestamp: Date.now()
      };

      setHistory(prev => prev.map(h => 
        h.id === entry.id 
          ? { ...h, response: errorResponse, status: 'error' }
          : h
      ));
    }
  };

  const handleKeyDown = (event: React.KeyboardEvent) => {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault();
      handleExecuteCommand();
    } else if (event.key === 'ArrowUp') {
      event.preventDefault();
      if (historyIndex < commandHistory.length - 1) {
        const newIndex = historyIndex + 1;
        setHistoryIndex(newIndex);
        setCommand(commandHistory[newIndex]);
      }
    } else if (event.key === 'ArrowDown') {
      event.preventDefault();
      if (historyIndex > 0) {
        const newIndex = historyIndex - 1;
        setHistoryIndex(newIndex);
        setCommand(commandHistory[newIndex]);
      } else if (historyIndex === 0) {
        setHistoryIndex(-1);
        setCommand('');
      }
    }
  };

  const handleClear = () => {
    setHistory([]);
  };

  const formatResponse = (response: RedisResponse): string => {
    switch (response.type) {
      case 'array':
        if (Array.isArray(response.value)) {
          return response.value.map((item, index) => `${index + 1}) ${item}`).join('\n');
        }
        return String(response.value);
      case 'null':
        return '(nil)';
      case 'error':
        return `ERROR: ${response.value}`;
      default:
        return String(response.value);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'success': return 'success';
      case 'error': return 'error';
      case 'pending': return 'warning';
      default: return 'default';
    }
  };

  return (
    <Paper sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <Toolbar variant="dense" sx={{ minHeight: '48px', borderBottom: '1px solid #333' }}>
        <Terminal fontSize="small" sx={{ mr: 1 }} />
        <Typography variant="h6" sx={{ flexGrow: 1, fontSize: '14px' }}>
          Redis Console
        </Typography>
        
        <Chip 
          label={connected ? 'Connected' : 'Disconnected'} 
          color={connected ? 'success' : 'error'} 
          size="small"
          sx={{ mr: 1 }}
        />
        
        <IconButton size="small" onClick={handleClear} disabled={history.length === 0}>
          <Clear fontSize="small" />
        </IconButton>
      </Toolbar>

      <Box 
        ref={listRef}
        sx={{ 
          flexGrow: 1, 
          overflow: 'auto',
          backgroundColor: '#0a0a0a',
          fontFamily: 'monospace'
        }}
      >
        <List dense sx={{ py: 0 }}>
          {history.map((entry) => (
            <React.Fragment key={entry.id}>
              <ListItem sx={{ py: 0.5, alignItems: 'flex-start' }}>
                <ListItemText
                  primary={
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Typography 
                        component="span" 
                        sx={{ 
                          fontFamily: 'monospace', 
                          fontSize: '13px',
                          color: '#2196f3'
                        }}
                      >
                        &gt; {entry.command}
                      </Typography>
                      <Chip 
                        label={entry.status} 
                        color={getStatusColor(entry.status) as any}
                        size="small"
                        sx={{ height: '16px', fontSize: '10px' }}
                      />
                    </Box>
                  }
                  secondary={
                    entry.response && (
                      <Typography 
                        component="pre" 
                        sx={{ 
                          fontFamily: 'monospace', 
                          fontSize: '12px',
                          color: entry.response.type === 'error' ? '#f44336' : '#e0e0e0',
                          whiteSpace: 'pre-wrap',
                          mt: 0.5,
                          mb: 0
                        }}
                      >
                        {formatResponse(entry.response)}
                      </Typography>
                    )
                  }
                />
              </ListItem>
            </React.Fragment>
          ))}
        </List>
      </Box>

      <Box sx={{ p: 1, borderTop: '1px solid #333' }}>
        <Box sx={{ display: 'flex', gap: 1 }}>
          <TextField
            fullWidth
            size="small"
            placeholder={connected ? "Enter Redis command (e.g., GRAPH.LIST, NODE.LIST graph1)" : "Not connected"}
            value={command}
            onChange={(e) => setCommand(e.target.value)}
            onKeyDown={handleKeyDown}
            disabled={!connected}
            sx={{
              '& .MuiInputBase-input': {
                fontFamily: 'monospace',
                fontSize: '13px'
              }
            }}
          />
          <IconButton 
            onClick={handleExecuteCommand}
            disabled={!command.trim() || !connected}
            color="primary"
          >
            <Send fontSize="small" />
          </IconButton>
        </Box>
        
        <Typography 
          variant="caption" 
          sx={{ 
            display: 'block', 
            mt: 0.5, 
            color: 'text.secondary',
            fontSize: '11px'
          }}
        >
          Press Enter to execute • Use ↑/↓ for command history
        </Typography>
      </Box>
    </Paper>
  );
};

export default RedisConsole;
