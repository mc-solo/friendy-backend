package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primarykey;default:gen_random_uuid()" json:"id"`
	Email        string     `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Username     string     `gorm:"size:25;column:username;uniqueIndex" json:"username"`
	PasswordHash string     `gorm:"column:password_hash;not null" json:"-"`
	FirstName    string     `gorm:"column:first_name;size:100" json:"first_name"`
	LastName     string     `gorm:"column:last_name;size:100" json:"last_name"`
	CreatedAt    *time.Time `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at;default:now()" json:"-"`
	DeletedAt    *time.Time `gorm:"column:deleted_at;index" json:"-"`

	// relations
	Preferences *UserPreference `gorm:"foreignKey:UserID" json:"preferences,omitempty"`
	Profile     *UserProfile    `gorm:"foreignKey:UserID" json:"user_profile,omitempty"`
}

func (User) TableName() string {
	return "users"
}
