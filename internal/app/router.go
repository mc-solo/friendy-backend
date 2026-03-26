package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"database/sql"

	httpDelivery "github.com/mc-solo/friendy/internal/delivery/http"
)

func NewRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// health check
	healthHandler := httpDelivery.NewHealthHandler(db)
	r.Get("/health", healthHandler.Check)

	return r
}
