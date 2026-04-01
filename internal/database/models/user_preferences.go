package models

import (
	"time"

	"github.com/google/uuid"
)

type UserPreference struct {
	ID        uuid.UUID  `gorm:"type:uuid;primarykey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"` // foreign key to users
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}
