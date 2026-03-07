package main

import (
	"log"

	"github.com/mc-solo/friendy/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: %w", err)
	}

	log.Printf("Config loaded: %+v", cfg)
}
