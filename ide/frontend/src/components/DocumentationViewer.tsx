import React, { useState, useEffect, useMemo } from 'react';
import {
  Box,
  Typography,
  TextField,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Paper,
  Divider,
  InputAdornment,
  Chip,
  Alert
} from '@mui/material';
import { Search, Description } from '@mui/icons-material';
import ReactMarkdown from 'react-markdown';
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { vscDarkPlus } from 'react-syntax-highlighter/dist/cjs/styles/prism';

interface DocFile {
  name: string;
  path: string;
  content: string;
}

interface DocumentationViewerProps {
  onClose?: () => void;
}

const DocumentationViewer: React.FC<DocumentationViewerProps> = ({ onClose }) => {
  const [docs, setDocs] = useState<DocFile[]>([]);
  const [selectedDoc, setSelectedDoc] = useState<DocFile | null>(null);
  const [searchTerm, setSearchTerm] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchDocs = async () => {
      try {
        const response = await fetch('http://localhost:8081/docs');
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        const docsData = await response.json();
        setDocs(docsData);
        if (docsData.length > 0) {
          setSelectedDoc(docsData[0]);
        }
      } catch (err) {
        setError('Failed to load documentation');
        console.error('Error fetching docs:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchDocs();
  }, []);

  const filteredDocs = useMemo(() => {
    if (!searchTerm) return docs;
    
    return docs.filter(doc => 
      doc.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
      doc.content.toLowerCase().includes(searchTerm.toLowerCase())
    );
  }, [docs, searchTerm]);

  const highlightedContent = useMemo(() => {
    if (!selectedDoc || !searchTerm) return selectedDoc?.content || '';
    
    const regex = new RegExp(`(${searchTerm})`, 'gi');
    return selectedDoc.content.replace(regex, '**$1**');
  }, [selectedDoc, searchTerm]);

  const searchMatches = useMemo(() => {
    if (!selectedDoc || !searchTerm) return 0;
    
    const regex = new RegExp(searchTerm, 'gi');
    const matches = selectedDoc.content.match(regex);
    return matches ? matches.length : 0;
  }, [selectedDoc, searchTerm]);

  if (loading) {
    return (
      <Box sx={{ p: 3, textAlign: 'center' }}>
        <Typography>Loading documentation...</Typography>
      </Box>
    );
  }

  if (error) {
    return (
      <Box sx={{ p: 3 }}>
        <Alert severity="error">{error}</Alert>
      </Box>
    );
  }

  return (
    <Box sx={{ height: '100vh', display: 'flex', flexDirection: 'column' }}>
      {/* Header */}
      <Box sx={{ 
        p: 2, 
        borderBottom: '1px solid #333',
        backgroundColor: '#1e1e1e'
      }}>
        <Typography variant="h6" sx={{ mb: 2, display: 'flex', alignItems: 'center' }}>
          <Description sx={{ mr: 1 }} />
          PathwayDB Documentation
        </Typography>
        
        <TextField
          fullWidth
          size="small"
          placeholder="Search documentation..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <Search />
              </InputAdornment>
            ),
          }}
          sx={{
            '& .MuiOutlinedInput-root': {
              backgroundColor: '#2d2d2d',
              '& fieldset': {
                borderColor: '#555',
              },
              '&:hover fieldset': {
                borderColor: '#777',
              },
              '&.Mui-focused fieldset': {
                borderColor: '#1976d2',
              },
            },
          }}
        />
        
        {searchTerm && searchMatches > 0 && (
          <Box sx={{ mt: 1 }}>
            <Chip 
              size="small" 
              label={`${searchMatches} matches found`}
              color="primary"
              variant="outlined"
            />
          </Box>
        )}
      </Box>

      {/* Main Content */}
      <Box sx={{ flexGrow: 1, display: 'flex', overflow: 'hidden' }}>
        {/* Sidebar - Document List */}
        <Box sx={{ 
          width: '300px', 
          borderRight: '1px solid #333',
          backgroundColor: '#1a1a1a',
          overflow: 'auto'
        }}>
          <List dense>
            {filteredDocs.map((doc) => (
              <ListItem key={doc.path} disablePadding>
                <ListItemButton
                  selected={selectedDoc?.path === doc.path}
                  onClick={() => setSelectedDoc(doc)}
                  sx={{
                    '&.Mui-selected': {
                      backgroundColor: '#2d2d2d',
                      '&:hover': {
                        backgroundColor: '#333',
                      },
                    },
                  }}
                >
                  <ListItemText 
                    primary={doc.name}
                    secondary={doc.path}
                    primaryTypographyProps={{
                      fontSize: '14px',
                      fontWeight: selectedDoc?.path === doc.path ? 'bold' : 'normal'
                    }}
                    secondaryTypographyProps={{
                      fontSize: '12px',
                      color: '#888'
                    }}
                  />
                </ListItemButton>
              </ListItem>
            ))}
          </List>
        </Box>

        {/* Main Content Area */}
        <Box sx={{ 
          flexGrow: 1, 
          overflow: 'auto',
          backgroundColor: '#1e1e1e'
        }}>
          {selectedDoc ? (
            <Box sx={{ p: 3 }}>
              <Typography variant="h5" sx={{ mb: 2, color: '#fff' }}>
                {selectedDoc.name}
              </Typography>
              <Divider sx={{ mb: 3, borderColor: '#333' }} />
              
              <Paper sx={{ 
                p: 3, 
                backgroundColor: '#2d2d2d',
                '& .markdown-content': {
                  color: '#e0e0e0',
                  lineHeight: 1.6,
                  '& h1, & h2, & h3, & h4, & h5, & h6': {
                    color: '#fff',
                    marginTop: '1.5em',
                    marginBottom: '0.5em',
                  },
                  '& h1': { fontSize: '2em', borderBottom: '1px solid #555', paddingBottom: '0.3em' },
                  '& h2': { fontSize: '1.5em', borderBottom: '1px solid #444', paddingBottom: '0.3em' },
                  '& h3': { fontSize: '1.25em' },
                  '& code': {
                    backgroundColor: '#1a1a1a',
                    padding: '2px 4px',
                    borderRadius: '3px',
                    fontSize: '0.9em',
                    color: '#f8f8f2',
                  },
                  '& pre': {
                    backgroundColor: '#1a1a1a',
                    padding: '1em',
                    borderRadius: '5px',
                    overflow: 'auto',
                  },
                  '& blockquote': {
                    borderLeft: '4px solid #555',
                    paddingLeft: '1em',
                    margin: '1em 0',
                    color: '#ccc',
                  },
                  '& table': {
                    borderCollapse: 'collapse',
                    width: '100%',
                    marginTop: '1em',
                    marginBottom: '1em',
                  },
                  '& th, & td': {
                    border: '1px solid #555',
                    padding: '8px 12px',
                    textAlign: 'left',
                  },
                  '& th': {
                    backgroundColor: '#333',
                    fontWeight: 'bold',
                  },
                  '& ul, & ol': {
                    paddingLeft: '2em',
                  },
                  '& li': {
                    marginBottom: '0.25em',
                  },
                  '& a': {
                    color: '#64b5f6',
                    textDecoration: 'none',
                    '&:hover': {
                      textDecoration: 'underline',
                    },
                  },
                  '& strong': {
                    backgroundColor: searchTerm ? '#ffeb3b' : 'transparent',
                    color: searchTerm ? '#000' : '#fff',
                    padding: searchTerm ? '1px 2px' : '0',
                    borderRadius: searchTerm ? '2px' : '0',
                  },
                }
              }}>
                <ReactMarkdown
                  components={{
                    code(props) {
                      const { children, className, ...rest } = props;
                      const match = /language-(\w+)/.exec(className || '');
                      return match ? (
                        <SyntaxHighlighter
                          style={vscDarkPlus}
                          language={match[1]}
                          PreTag="div"
                        >
                          {String(children).replace(/\n$/, '')}
                        </SyntaxHighlighter>
                      ) : (
                        <code className={className} {...rest}>
                          {children}
                        </code>
                      );
                    },
                  }}
                >
                  {highlightedContent}
                </ReactMarkdown>
              </Paper>
            </Box>
          ) : (
            <Box sx={{ 
              p: 3, 
              textAlign: 'center',
              color: '#888'
            }}>
              <Typography>Select a document to view</Typography>
            </Box>
          )}
        </Box>
      </Box>
    </Box>
  );
};

export default DocumentationViewer;
