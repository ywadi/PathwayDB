package redis

import "time"

// Config holds the configuration for the Redis server
type Config struct {
	// Server address to bind to
	Address string
	
	// Maximum number of concurrent connections
	MaxConnections int
	
	// Connection timeout
	ConnectionTimeout time.Duration
	
	// Read timeout for commands
	ReadTimeout time.Duration
	
	// Write timeout for responses
	WriteTimeout time.Duration
	
	// Enable debug logging
	Debug bool
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Address:           ":6379",
		MaxConnections:    1000,
		ConnectionTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		Debug:             false,
	}
}
