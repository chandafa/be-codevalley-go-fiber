package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Conditions map[string]interface{}

func (c Conditions) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Conditions) Scan(value interface{}) error {
	if value == nil {
		*c = make(Conditions)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, c)
}

type Achievement struct {
	ID          uuid.UUID  `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	Title       string     `json:"title" gorm:"not null" validate:"required"`
	Description string     `json:"description" gorm:"type:text" validate:"required"`
	RewardCoins int        `json:"reward_coins" gorm:"default:0"`
	RewardEXP   int        `json:"reward_exp" gorm:"default:0"`
	Conditions  Conditions `json:"conditions" gorm:"type:json"`
	IconURL     string     `json:"icon_url"`
	IsActive    bool       `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relationships
	UserAchievements []UserAchievement `json:"user_achievements,omitempty" gorm:"foreignKey:AchievementID"`
}

func (a *Achievement) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

type UserAchievement struct {
	ID            uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID        uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index"`
	AchievementID uuid.UUID `json:"achievement_id" gorm:"type:char(36);not null;index"`
	UnlockedAt    time.Time `json:"unlocked_at"`

	// Relationships
	User        User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Achievement Achievement `json:"achievement,omitempty" gorm:"foreignKey:AchievementID"`
}

func (ua *UserAchievement) BeforeCreate(tx *gorm.DB) error {
	if ua.ID == uuid.Nil {
		ua.ID = uuid.New()
	}
	return nil
}