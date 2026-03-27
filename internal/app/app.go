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

type App struct {
	Router   *chi.Mux
	Config   *config.Config
	DB       *gorm.DB
	Handlers *Handlers // defined in handlers.go
}

func New(cfg *config.Config) (*App, error) {
	// open database
	db, err := cfg.OpenDB()
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// gets *sql.DB for health check
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from gorm: %w", err)
	}

	userStorePtr := store.NewUserStore(db)
	refreshStorePtr := store.NewRefreshTokenStore(db)

	// token config
	tokenCfg := token.Config{
		AccessSecret:  cfg.JWT.Secret,
		RefreshSecret: cfg.JWT.Secret,
		AccessExpiry:  cfg.JWT.AccessTokenExp,
		RefreshExpiry: cfg.JWT.RefreshTokenExp,
	}

	// auth service
	authSvc := authService.NewService(*userStorePtr, *refreshStorePtr, tokenCfg)

	// handlers
	healthHandler := health.NewHandler(sqlDB)
	authHandler := auth.NewHandler(authSvc)

	handlers := &Handlers{
		Health: healthHandler,
		Auth:   authHandler,
		// add handlers here [like chat, verification...]
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

func (a *App) Run() error {
	addr := fmt.Sprintf(":%d", a.Config.ServerPort)
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, a.Router)
}

func (a *App) Close() error {
	sqlDB, err := a.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
