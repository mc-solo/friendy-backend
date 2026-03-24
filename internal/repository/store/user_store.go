package store

import (
	"context"
	"errors"

	"github.com/mc-solo/friendy/internal/database/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// defines the data access interface for users
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (user *models.User)
	GetByID(ctx context.Context, id uuid.UUID) (user *models.User)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// implements UserRepository with gorm
type UserStore struct {
	db *gorm.DB
}

// creates a new user store with the given db conn
func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db: db}
}

// inserts a new user into db
func (s *UserStore) Create(ctx context.Context, user *models.User) error {
	return s.db.WithContext(ctx).Create(user).Error
}

// retrieves a user by email [returns nil, nil if not found]
func (s *UserStore) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).Where("email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //cuz this is not our problem
		}
		return nil, err
	}
	return &user, nil
}

// retrieves a user by id
func (s *UserStore) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := s.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// updates an existing user. the user must have a valid ID
func (s *UserStore) Update(ctx context.Context, user *models.User) error {
	return s.db.WithContext(ctx).Save(user).Error
}

func (s *UserStore) Delete(ctx context.Context, id uuid.UUID) error {
	result := s.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	// i should prolly check RowsAffected here
	return nil
}
