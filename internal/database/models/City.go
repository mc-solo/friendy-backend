package database

import (
	uuid "github.com/google/uuid"
	"time"
)

type City struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
    Name      string    `gorm:"not null;uniqueIndex:idx_city_country_name"`
    CountryID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_city_country_name"`
	Country Country `gorm:"foreignKey:CountryID"`

    Latitude  float64
    Longitude float64

    CreatedAt time.Time
    UpdatedAt time.Time
}

func (City) TableName() string {
	return "cities"
}