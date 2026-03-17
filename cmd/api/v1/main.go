package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	// open db conn
	db, err := cfg.OpenDB()
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	log.Println("Databse connected succesfully")

	// health check
	healthHandler := func(w http.ResponseWriter, r *http.Request) {
		sqlDB, err := db.DB()
		if err != nil {
			http.Error(w, "Cannot get underlying DB", http.StatusInternalServerError)
			return
		}

		if err := sqlDB.Ping(); err != nil {
			http.Error(w, "Databse ping failed", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("healthy"))

	}

	// setup router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", healthHandler)

	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}

}
