package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StoryProgress struct {
	ID          uuid.UUID  `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:char(36);not null;index"`
	Chapter     int        `json:"chapter" gorm:"not null"`
	Milestone   string     `json:"milestone" gorm:"not null" validate:"required"`
	IsCompleted bool       `json:"is_completed" gorm:"default:false"`
	UnlockedAt  *time.Time `json:"unlocked_at"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (sp *StoryProgress) BeforeCreate(tx *gorm.DB) error {
	if sp.ID == uuid.Nil {
		sp.ID = uuid.New()
	}
	return nil
}