package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sainudheenp/goecom/config"
	"github.com/sainudheenp/goecom/server"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting e-commerce server in %s mode", cfg.Server.Env)

	// Create server
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer srv.Close()

	// Handle graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")
		srv.Close()
		os.Exit(0)
	}()

	// Run server
	if err := srv.Run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
