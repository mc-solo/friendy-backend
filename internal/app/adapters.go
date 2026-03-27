package app

import (
	"github.com/google/uuid"
	"github.com/mc-solo/friendy/internal/config"
	"github.com/mc-solo/friendy/internal/utils/password"
	"github.com/mc-solo/friendy/internal/utils/token"
)

// adapts the password pack to the expected interface
type passwordHasher struct{}

func (p passwordHasher) Hash(pwd string) (string, error) {
	return password.Hash(pwd)
}

func (p passwordHasher) Check(pwd, hash string) bool {
	return password.Check(pwd, hash)
}

// adapts the token pack to the interface
type tokenMaker struct {
	cfg token.Config
}

func NewTokenMaker(cfg token.Config) *tokenMaker {
	return &tokenMaker{cfg: cfg}
}

func (t *tokenMaker) GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
	return token.GenAccessToken(userID, email, t.cfg)
}

func (t *tokenMaker) GenerateRefreshToken(userID uuid.UUID) (string, error) {
	return token.GenRefreshToken(userID, t.cfg)
}

func (t *tokenMaker) ValidateAccessToken(tokenString string) (*token.Claims, error) {
	return token.ValidateAccessToken(tokenString, t.cfg)
}

func (t *tokenMaker) ValidateRefreshToken(tokenString string) (*token.Claims, error) {
	return token.ValidateRefreshToken(tokenString, t.cfg)
}
