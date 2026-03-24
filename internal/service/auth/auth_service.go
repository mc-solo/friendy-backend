package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/mc-solo/friendy/internal/database/models"
	"github.com/mc-solo/friendy/internal/repository/store"
	"github.com/mc-solo/friendy/internal/utils/password"
	"github.com/mc-solo/friendy/internal/utils/token"
)

type Service struct {
	userStore store.UserStore
	tokenCfg  token.Config
}

func NewService(userStore store.UserStore, tokenCfg token.Config) *Service {
	return &Service{
		userStore: userStore,
		tokenCfg:  tokenCfg,
	}
}

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already registered")
)

func (s *Service) Register(ctx context.Context, email, plainPassword string) (*models.User, error) {

	// check if the user exists
	existing, _ := s.userStore.GetByEmail(ctx, email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	// hash the password
	hashed, err := password.Hash(plainPassword)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}

	user := &models.User{
		Email:        email,
		PasswordHash: hashed,
	}

	if err := s.userStore.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("creating user: %w", err)
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, email, plainPassword string) (accessToken, refreshToken string, err error) {
	user, err := s.userStore.GetByEmail(ctx, email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	if !password.Check(plainPassword, user.PasswordHash) {
		return "", "", ErrInvalidCredentials
	}

	// then access token
	access, err := token.GenAccessToken(user.ID, user.Email, s.tokenCfg)
	if err != nil {
		return "", "", fmt.Errorf("generating refresh token: %w", err)
	}

	// and here the refresh token
	refresh, err := token.GenRefreshToken(user.ID, s.tokenCfg)
	if err != nil {
		return "", "", fmt.Errorf("generating refresh token: %w", err)
	}

	// TODO: store refresh token hash in DB

	return access, refresh, nil
}

// TODO: i'll implement logout here once i'm done with my login/register perfectly
func (s *Service) Refresh(ctx context.Context, refreshTokenStr string) (newAccessToken string, err error) {
	claims, err := token.ValidateRefreshToken(refreshTokenStr, s.tokenCfg)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}
	// TODO: verify the refresh token matches the stored hash

	return token.GenAccessToken(claims.UserID, claims.Email, s.tokenCfg)
}
