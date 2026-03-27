package app

import (
	"github.com/mc-solo/friendy/internal/delivery/http/auth"
	"github.com/mc-solo/friendy/internal/delivery/http/health"
)

type Handlers struct {
	Health *health.Handler
	Auth   *auth.Handler
}
