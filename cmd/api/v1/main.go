package main

import (
	"log"

	"github.com/mc-solo/friendy/internal/app"
	"github.com/mc-solo/friendy/internal/config"
)

func main() {
	// load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// validate req fields
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	application := app.New(cfg)
	if err := application.Run(); err != nil {
		log.Fatalf("app failed: %v", err)
	}
}
