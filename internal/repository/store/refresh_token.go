package store

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/mc-solo/friendy/internal/database/models"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, token *models.RefreshToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
	Delete(ctx context.Context, id uuid.UUID) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
}

type RefreshTokenStore struct {
	db *gorm.DB
}

func NewRefreshTokenStore(db *gorm.DB) *RefreshTokenStore {
	return &RefreshTokenStore{db: db}
}

func (s *RefreshTokenStore) Create(ctx context.Context, token *models.RefreshToken) error {
	return s.db.WithContext(ctx).Create(token).Error
}

func (s *RefreshTokenStore) GetByTokenHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := s.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (s *RefreshTokenStore) Delete(ctx context.Context, id uuid.UUID) error {
	return s.db.WithContext(ctx).Delete(&models.RefreshToken{}, "id = ?", id).Error
}

func (s *RefreshTokenStore) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	return s.db.WithContext(ctx).Delete(&models.RefreshToken{}, "user_id = ?", userID).Error
}
