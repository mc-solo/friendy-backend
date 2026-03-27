package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/mc-solo/internal/service/auth"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	// global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// the v1 api
	apiRouter := r.Route("/api/v1", func(api chi.Router) {
		authService := auth.NewService()
		authHandler := auth.NewHandler(authService)

		// mount auth routes
		api.Route("/auth", func(authRouter chi.Router) {
			auth.RegisterRoutes(authRouter, authHandler)
		})
	})

	return r
}
