package database

import (
	uuid "github.com/google/uuid"
	"time"
)


type Country struct {
    ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Name string    `gorm:"unique;not null"`

    CreatedAt time.Time
    UpdatedAt time.Time
}

func (Country) TableName() string {
	return "countries"
}
