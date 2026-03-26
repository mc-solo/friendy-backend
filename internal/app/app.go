package app

import (
	"fmt"
	"net/http"

	"github.com/mc-solo/friendy/internal/config"
)

type App struct {
	cfg *config.Config
}

func New(cfg *config.Config) *App {
	return &App{cfg: cfg}
}

func (a *App) Run() error {
	// db
	db, err := a.cfg.OpenDB()
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// router
	router := NewRouter(sqlDB)

	addr := fmt.Sprintf(":%d", a.cfg.ServerPort)

	return http.ListenAndServe(addr, router)
}
