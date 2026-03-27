package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mc-solo/friendy/internal/app"
	"github.com/mc-solo/friendy/internal/config"
)

func main() {
	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// init app
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialise app: %v", err)
	}
	defer func() {
		if err := application.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// create HTTP server with the router
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	server := &http.Server{
		Addr:    addr,
		Handler: application.Router,
	}

	// start server in goroutine
	go func() {
		log.Printf("Starting server on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}
