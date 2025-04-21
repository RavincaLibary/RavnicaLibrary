package main

import (
	"RavnicaLibrary/internal/api"
	"log"
	"os"
)

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create and start the server
	server := api.NewServer(port)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
