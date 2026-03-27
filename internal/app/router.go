package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/mc-solo/friendy/internal/delivery/http/auth"
	"github.com/mc-solo/friendy/internal/delivery/http/health"
)

// NewRouter creates the main router with all endpoints.
func NewRouter(h *Handlers) *chi.Mux {
	r := chi.NewRouter()

	// gloabl middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)

	// health check
	health.RegisterRoutes(r, h.Health)

	// API v1 group
	r.Route("/api/v1", func(api chi.Router) {
		// auth routes under /api/v1/auth
		api.Route("/auth", func(authRouter chi.Router) {
			auth.RegisterRoutes(authRouter, h.Auth)
		})

		// other modules will be added here:
		// api.Route("/users", user.RegisterRoutes)
		// api.Route("/chat", chat.RegisterRoutes)
	})

	return r
}
