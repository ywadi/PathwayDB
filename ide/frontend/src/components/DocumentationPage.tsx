import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Paper, 
  Typography, 
  List, 
  ListItem, 
  ListItemButton, 
  ListItemText, 
  CircularProgress, 
  IconButton,
  Chip,
  Fade,
  useTheme
} from '@mui/material';
import { 
  Close, 
  Description, 
  Code, 
  Settings, 
  BugReport,
  Article
} from '@mui/icons-material';

interface DocumentationPageProps {
  onClose: () => void;
}

const getFileIcon = (filename: string) => {
  if (filename === 'README.md') return <Article sx={{ fontSize: 18 }} />;
  if (filename === 'COMMANDS.md') return <Code sx={{ fontSize: 18 }} />;
  if (filename === 'IDE.md') return <Settings sx={{ fontSize: 18 }} />;
  if (filename === 'TESTS.md') return <BugReport sx={{ fontSize: 18 }} />;
  return <Description sx={{ fontSize: 18 }} />;
};

const getFileDisplayName = (filename: string) => {
  const names: { [key: string]: string } = {
    'README.md': 'Project Overview',
    'COMMANDS.md': 'Redis Commands',
    'IDE.md': 'IDE Guide',
    'TESTS.md': 'Test Suite'
  };
  return names[filename] || filename;
};

const DocumentationPage: React.FC<DocumentationPageProps> = ({ onClose }) => {
  const theme = useTheme();
  const [docFiles, setDocFiles] = useState<string[]>([]);
  const [selectedFile, setSelectedFile] = useState<string>('');
  const [content, setContent] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchDocFiles = async () => {
      try {
        const response = await fetch('http://localhost:8081/api/docs');
        const files = await response.json();
        setDocFiles(files);
        if (files.length > 0) {
          handleFileSelect(files[0]);
        }
      } catch (error) {
        console.error('Failed to fetch doc files:', error);
        setContent('<div style="color: #ff6b6b; padding: 20px; text-align: center;"><h3>⚠️ Error loading documentation</h3><p>Please ensure the backend server is running.</p></div>');
      } finally {
        setLoading(false);
      }
    };

    fetchDocFiles();
  }, []);

  const handleFileSelect = async (filename: string) => {
    setSelectedFile(filename);
    setLoading(true);
    try {
      const response = await fetch(`http://localhost:8081/api/docs/${filename}`);
      const html = await response.text();
      setContent(html);
    } catch (error) {
      console.error(`Failed to fetch ${filename}:`, error);
      setContent(`<div style="color: #ff6b6b; padding: 20px; text-align: center;"><h3>⚠️ Error loading ${filename}</h3><p>Failed to fetch document content.</p></div>`);
    } finally {
      setLoading(false);
    }
  };

  const handleClose = () => {
    // Clear any potential references before closing
    setContent('');
    setSelectedFile('');
    setDocFiles([]);
    // Use setTimeout to ensure state is cleared before calling onClose
    setTimeout(() => {
      onClose();
    }, 0);
  };

  return (
    <Fade in timeout={300}>
      <Box sx={{ 
        display: 'flex', 
        height: '100vh', 
        backgroundColor: '#0d1117',
        color: '#e6edf3',
        fontFamily: '"Segoe UI", "Helvetica Neue", Arial, sans-serif'
      }}>
        {/* Sidebar */}
        <Paper elevation={0} sx={{ 
          width: '320px', 
          height: '100%', 
          overflowY: 'auto',
          backgroundColor: '#161b22',
          borderRight: '1px solid #30363d',
          borderRadius: 0
        }}>
          {/* Header */}
          <Box sx={{ 
            p: 2, 
            borderBottom: '1px solid #30363d', 
            display: 'flex', 
            justifyContent: 'space-between', 
            alignItems: 'center',
            background: 'linear-gradient(135deg, #1f2937 0%, #161b22 100%)'
          }}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <Article sx={{ color: '#58a6ff', fontSize: 20 }} />
              <Typography variant="h6" sx={{ 
                fontSize: '16px', 
                fontWeight: 600,
                color: '#e6edf3'
              }}>
                Documentation
              </Typography>
            </Box>
            <IconButton 
              onClick={handleClose} 
              sx={{ 
                color: '#8b949e',
                '&:hover': {
                  backgroundColor: '#30363d',
                  color: '#e6edf3'
                }
              }}
            >
              <Close />
            </IconButton>
          </Box>

          {/* File List */}
          <List sx={{ p: 1 }}>
            {docFiles.map((file) => (
              <ListItem key={file} disablePadding sx={{ mb: 0.5 }}>
                <ListItemButton 
                  selected={selectedFile === file}
                  onClick={() => handleFileSelect(file)}
                  sx={{
                    borderRadius: '6px',
                    mx: 1,
                    '&.Mui-selected': {
                      backgroundColor: '#1f6feb',
                      color: 'white',
                      '&:hover': {
                        backgroundColor: '#1f6feb',
                      }
                    },
                    '&:hover': {
                      backgroundColor: '#30363d',
                    }
                  }}
                >
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5, width: '100%' }}>
                    {getFileIcon(file)}
                    <Box sx={{ flexGrow: 1 }}>
                      <Typography variant="body2" sx={{ 
                        fontWeight: selectedFile === file ? 600 : 400,
                        color: selectedFile === file ? 'white' : '#e6edf3'
                      }}>
                        {getFileDisplayName(file)}
                      </Typography>
                      <Typography variant="caption" sx={{ 
                        color: selectedFile === file ? 'rgba(255,255,255,0.7)' : '#8b949e',
                        fontSize: '11px'
                      }}>
                        {file}
                      </Typography>
                    </Box>
                  </Box>
                </ListItemButton>
              </ListItem>
            ))}
          </List>

          {/* Footer */}
          <Box sx={{ p: 2, mt: 'auto', borderTop: '1px solid #30363d' }}>
            <Chip 
              label={`${docFiles.length} documents`}
              size="small"
              sx={{ 
                backgroundColor: '#30363d',
                color: '#8b949e',
                fontSize: '11px'
              }}
            />
          </Box>
        </Paper>

        {/* Content Area */}
        <Box sx={{ 
          flexGrow: 1, 
          display: 'flex',
          flexDirection: 'column',
          height: '100vh'
        }}>
          {/* Content Header */}
          {selectedFile && (
            <Box sx={{ 
              p: 2, 
              borderBottom: '1px solid #30363d',
              backgroundColor: '#0d1117'
            }}>
              <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                {getFileIcon(selectedFile)}
                <Typography variant="h5" sx={{ 
                  fontWeight: 600,
                  color: '#e6edf3'
                }}>
                  {getFileDisplayName(selectedFile)}
                </Typography>
                <Chip 
                  label={selectedFile}
                  size="small"
                  sx={{ 
                    backgroundColor: '#30363d',
                    color: '#8b949e',
                    fontSize: '11px'
                  }}
                />
              </Box>
            </Box>
          )}

          {/* Content */}
          <Box sx={{ 
            flexGrow: 1, 
            overflowY: 'auto',
            backgroundColor: '#0d1117'
          }}>
            {loading ? (
              <Box sx={{ 
                display: 'flex', 
                justifyContent: 'center', 
                alignItems: 'center', 
                height: '100%',
                flexDirection: 'column',
                gap: 2
              }}>
                <CircularProgress sx={{ color: '#58a6ff' }} />
                <Typography sx={{ color: '#8b949e' }}>Loading documentation...</Typography>
              </Box>
            ) : (
              <Box
                className="markdown-body"
                dangerouslySetInnerHTML={{ __html: content }}
                sx={{
                  p: 4,
                  maxWidth: '900px',
                  margin: '0 auto',
                  lineHeight: 1.6,
                  
                  // Typography
                  '& h1': {
                    color: '#e6edf3',
                    fontSize: '2rem',
                    fontWeight: 600,
                    borderBottom: '1px solid #30363d',
                    paddingBottom: '0.5em',
                    marginBottom: '1em',
                    marginTop: '1.5em'
                  },
                  '& h2': {
                    color: '#e6edf3',
                    fontSize: '1.5rem',
                    fontWeight: 600,
                    borderBottom: '1px solid #30363d',
                    paddingBottom: '0.3em',
                    marginBottom: '0.8em',
                    marginTop: '1.2em'
                  },
                  '& h3': {
                    color: '#e6edf3',
                    fontSize: '1.25rem',
                    fontWeight: 600,
                    marginBottom: '0.6em',
                    marginTop: '1em'
                  },
                  '& h4, & h5, & h6': {
                    color: '#e6edf3',
                    fontWeight: 600,
                    marginBottom: '0.5em',
                    marginTop: '0.8em'
                  },
                  '& p': {
                    color: '#e6edf3',
                    marginBottom: '1em',
                    fontSize: '14px'
                  },
                  '& li': {
                    color: '#e6edf3',
                    marginBottom: '0.3em',
                    fontSize: '14px'
                  },
                  
                  // Links
                  '& a': {
                    color: '#58a6ff',
                    textDecoration: 'none',
                    '&:hover': {
                      textDecoration: 'underline'
                    }
                  },
                  
                  // Code blocks
                  '& pre': {
                    backgroundColor: '#161b22',
                    border: '1px solid #30363d',
                    borderRadius: '6px',
                    padding: '16px',
                    overflow: 'auto',
                    fontSize: '13px',
                    lineHeight: 1.45,
                    marginBottom: '1em'
                  },
                  '& code': {
                    fontFamily: '"SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace',
                    fontSize: '85%',
                    backgroundColor: '#161b22',
                    padding: '0.2em 0.4em',
                    borderRadius: '3px',
                    border: '1px solid #30363d'
                  },
                  '& pre code': {
                    backgroundColor: 'transparent',
                    border: 'none',
                    padding: 0,
                    fontSize: 'inherit'
                  },
                  
                  // Tables
                  '& table': {
                    borderCollapse: 'collapse',
                    width: '100%',
                    marginBottom: '1em',
                    border: '1px solid #30363d',
                    borderRadius: '6px',
                    overflow: 'hidden'
                  },
                  '& th': {
                    backgroundColor: '#161b22',
                    color: '#e6edf3',
                    fontWeight: 600,
                    padding: '8px 12px',
                    border: '1px solid #30363d',
                    textAlign: 'left'
                  },
                  '& td': {
                    color: '#e6edf3',
                    padding: '8px 12px',
                    border: '1px solid #30363d'
                  },
                  '& tr:nth-of-type(even)': {
                    backgroundColor: '#0d1117'
                  },
                  
                  // Lists
                  '& ul, & ol': {
                    paddingLeft: '2em',
                    marginBottom: '1em'
                  },
                  
                  // Blockquotes
                  '& blockquote': {
                    borderLeft: '4px solid #30363d',
                    paddingLeft: '1em',
                    marginLeft: 0,
                    color: '#8b949e',
                    fontStyle: 'italic'
                  },
                  
                  // Strong/Bold
                  '& strong, & b': {
                    color: '#e6edf3',
                    fontWeight: 600
                  }
                }}
              />
            )}
          </Box>
        </Box>
      </Box>
    </Fade>
  );
};

export default DocumentationPage;
