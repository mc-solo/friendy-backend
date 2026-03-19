package models

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	ID                    uuid.UUID        `gorm:"type:uuid;primarykey;default:gen_random_uuid()" json:"id"`
	UserID                uuid.UUID        `gorm:"type:uuid;uniqueIndex;not null;" json:"user_id"` //foreign key
	Gender                string           `gorm:"size:20" json:"gender"`
	Bio                   string           `gorm:"type:varchar(255);" json:"bio"`
	BodyType              BodyType         `gorm:"type:varchar(100)" json:"body_type"`
	BirthDate             *time.Time       `gorm:"default:"`
	HeightCm              float64          `gorm:"type:decimal(5,2)" json:"height_cm"`
	Language              Language         `gorm:"type:varchar(10);default:'en'" json:"lang"` //TODO: update this type so users can select multiple langs
	EduLevel              EducatoinalLevel `gorm:"column:edu_level;type:varchar(50)" json:"edu_level"`
	ProfileCompletionRate int              `gorm:"default:0" json:"profile_completion_rate"`
	// will add more fields as we go

	CreatedAt time.Time `gorm:"column:created_at;default:now()" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;default:now()" json:"updated_at"`
}
