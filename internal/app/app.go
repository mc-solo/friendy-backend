package app

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"gorm.io/gorm"

	"github.com/mc-solo/friendy/internal/config"
	"github.com/mc-solo/friendy/internal/delivery/http/auth"
	"github.com/mc-solo/friendy/internal/delivery/http/health"
	"github.com/mc-solo/friendy/internal/repository/store"
	authService "github.com/mc-solo/friendy/internal/service/auth"
	"github.com/mc-solo/friendy/internal/utils/token"
)

// app holds the core components of the application.
type App struct {
	Router   *chi.Mux
	Config   *config.Config
	DB       *gorm.DB
	Handlers *Handlers // defined in handlers.go
}

// New initialises the entire application.
func New(cfg *config.Config) (*App, error) {
	// open db conn using gorm
	db, err := cfg.OpenDB()
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// gets *sql.DB for health checks
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm: %w", err)
	}

	// repositories
	userStore := store.NewUserStore(db)
	refreshStore := store.NewRefreshTokenStore(db)
	// add repositories here when you create them later

	// utilities [ adapters from adapters.go ]
	hasher := passwordHasher{}

	tokenCfg := token.Config{
		AccessSecret:  cfg.JWT.Secret,
		RefreshSecret: cfg.JWT.Secret,
		AccessExpiry:  cfg.JWT.AccessTokenExp,
		RefreshExpiry: cfg.JWT.RefreshTokenExp,
	}
	tokenMaker := &tokenMaker{cfg: tokenCfg} // from adapters.go

	// services
	authSvc := authService.NewService(userStore, refreshStore, hasher, tokenMaker)

	// handlers
	healthHandler := health.NewHandler(sqlDB)
	authHandler := auth.NewHandler(authSvc)

	handlers := &Handlers{
		Health: healthHandler,
		Auth:   authHandler,
		// you guys can add handlers here as you go
	}

	// router
	router := NewRouter(handlers)

	return &App{
		Router:   router,
		Config:   cfg,
		DB:       db,
		Handlers: handlers,
	}, nil
}

// Run starts the HTTP server.
func (a *App) Run() error {
	addr := fmt.Sprintf(":%d", a.Config.ServerPort)
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, a.Router)
}

// Close gracefully closes the db conn
func (a *App) Close() error {
	sqlDB, err := a.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
