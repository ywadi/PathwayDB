package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Get the project root directory
	projectRoot, err := filepath.Abs("../../")
	if err != nil {
		log.Fatal("Failed to get project root:", err)
	}

	// Change to the IDE backend directory
	ideBackendDir := filepath.Join(projectRoot, "ide", "backend")
	if err := os.Chdir(ideBackendDir); err != nil {
		log.Fatal("Failed to change to IDE backend directory:", err)
	}

	// Execute the IDE server
	cmd := exec.Command("go", "run", "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	log.Println("Starting PathwayDB IDE Server...")
	if err := cmd.Run(); err != nil {
		log.Fatal("IDE server failed:", err)
	}
}
