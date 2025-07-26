package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BadgeType string

const (
	BadgeTypeAchievement BadgeType = "achievement"
	BadgeTypeEvent       BadgeType = "event"
	BadgeTypeSpecial     BadgeType = "special"
	BadgeTypeSeasonal    BadgeType = "seasonal"
)

type BadgeRarity string

const (
	BadgeRarityCommon    BadgeRarity = "common"
	BadgeRarityUncommon  BadgeRarity = "uncommon"
	BadgeRarityRare      BadgeRarity = "rare"
	BadgeRarityEpic      BadgeRarity = "epic"
	BadgeRarityLegendary BadgeRarity = "legendary"
)

type Badge struct {
	ID          uuid.UUID   `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	Name        string      `json:"name" gorm:"not null;index" validate:"required"`
	Description string      `json:"description" gorm:"type:text"`
	Type        BadgeType   `json:"type" gorm:"type:enum('achievement','event','special','seasonal');not null"`
	Rarity      BadgeRarity `json:"rarity" gorm:"type:enum('common','uncommon','rare','epic','legendary');default:'common'"`
	IconURL     string      `json:"icon_url"`
	Conditions  Conditions  `json:"conditions" gorm:"type:json"`
	IsActive    bool        `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`

	// Relationships
	UserBadges []UserBadge `json:"user_badges,omitempty" gorm:"foreignKey:BadgeID"`
}

func (b *Badge) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type UserBadge struct {
	ID       uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID   uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index"`
	BadgeID  uuid.UUID `json:"badge_id" gorm:"type:char(36);not null;index"`
	EarnedAt time.Time `json:"earned_at"`

	// Relationships
	User  User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Badge Badge `json:"badge,omitempty" gorm:"foreignKey:BadgeID"`
}

func (ub *UserBadge) BeforeCreate(tx *gorm.DB) error {
	if ub.ID == uuid.Nil {
		ub.ID = uuid.New()
	}
	return nil
}