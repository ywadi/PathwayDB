package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ywadi/PathwayDB/redis"
	"github.com/ywadi/PathwayDB/storage"
)

// getEnv reads an environment variable or returns a fallback value.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	// Configuration with environment variable override
	redisAddr := getEnv("REDIS_ADDR", ":6379")

	// Command line flags
	var (
		addr     = flag.String("addr", redisAddr, "Redis server address")
		dataDir  = flag.String("data", "./data", "Data directory for storage")
		debug    = flag.Bool("debug", false, "Enable debug logging")
	)
	flag.Parse()

	// Create storage engine
	storageEngine := storage.NewBadgerEngine()
	if err := storageEngine.Open(*dataDir); err != nil {
		log.Fatalf("Failed to open storage engine: %v", err)
	}
	defer storageEngine.Close()

	// Create Redis server configuration
	config := redis.DefaultConfig()
	config.Address = *addr
	config.Debug = *debug

	// Create and start Redis server
	server := redis.NewServer(config, storageEngine)

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down PathwayDB Redis server...")
		server.Stop()
		os.Exit(0)
	}()

	log.Printf("PathwayDB Redis server starting on %s", config.Address)
	log.Printf("Data directory: %s", *dataDir)
	
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
