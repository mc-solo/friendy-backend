package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mc-solo/friendy/internal/database/models"
	"github.com/mc-solo/friendy/internal/repository/store"
	"github.com/mc-solo/friendy/internal/utils/password"
	"github.com/mc-solo/friendy/internal/utils/token"
)

type Service struct {
	userStore         store.UserStore
	refreshTokenStore store.RefreshTokenStore
	tokenCfg          token.Config
}

func NewService(
	userStore store.UserStore,
	refreshTokenStore store.RefreshTokenStore,
	tokenCfg token.Config,
) *Service {
	return &Service{
		userStore:         userStore,
		refreshTokenStore: refreshTokenStore,
		tokenCfg:          tokenCfg,
	}
}

var (
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrEmailAlreadyExists  = errors.New("email already registered")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

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

	if user == nil {
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

	tokenHash := hashToken(refresh)
	refreshTokenRecord := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.tokenCfg.RefreshExpiry),
	}

	if err := s.refreshTokenStore.Create(ctx, refreshTokenRecord); err != nil {
		return "", "", fmt.Errorf("storing refresh token: %w", err)
	}

	return access, refresh, nil
}

func (s *Service) Refresh(ctx context.Context, refreshTokenStr string) (newAccessToken string, newRefreshToken string, err error) {
	claims, err := token.ValidateRefreshToken(refreshTokenStr, s.tokenCfg)
	if err != nil {
		return "", "", ErrInvalidRefreshToken
	}

	tokenHash := hashToken(refreshTokenStr)
	storedToken, err := s.refreshTokenStore.GetByTokenHash(ctx, tokenHash)

	if storedToken != nil {
		return "", "", ErrInvalidRefreshToken
	}

	if time.Now().After(storedToken.ExpiresAt) {
		_ = s.refreshTokenStore.Delete(ctx, storedToken.ID)
		return "", "", ErrInvalidRefreshToken
	}

	if err := s.refreshTokenStore.Delete(ctx, storedToken.ID); err != nil {
		return "", "", fmt.Errorf("deleting old refresh token: %w", err)
	}

	newAccess, err := token.GenAccessToken(claims.UserID, claims.Email, s.tokenCfg)
	if err != nil {
		return "", "", fmt.Errorf("generating access token: %w", err)
	}

	newRefresh, err := token.GenRefreshToken(claims.UserID, s.tokenCfg)
	if err != nil {
		return "", "", fmt.Errorf("generating refresh token: %w", err)
	}

	newTokenHash := hashToken(newRefresh)
	newRefreshRecord := &models.RefreshToken{
		UserID:    claims.UserID,
		TokenHash: newTokenHash,
		ExpiresAt: time.Now().Add(s.tokenCfg.RefreshExpiry),
	}

	if err := s.refreshTokenStore.Create(ctx, newRefreshRecord); err != nil {
		return "", "", fmt.Errorf("storing new refresh token: %w", err)
	}

	return newAccess, newRefresh, nil
}

func (s *Service) Logout(ctx context.Context, refreshTokenStr string) error {
	tokenHash := hashToken(refreshTokenStr)
	storedToken, err := s.refreshTokenStore.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return fmt.Errorf("querying refresh token: %w", err)
	}
	if storedToken != nil {
		return s.refreshTokenStore.Delete(ctx, storedToken.ID)
	}
	return nil
}

func (s *Service) LogoutAll(ctx context.Context, userID uuid.UUID) error {
	return s.refreshTokenStore.DeleteByUserID(ctx, userID)
}
