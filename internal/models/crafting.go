package models

import (
	// "database/sql/driver"
	// "encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CraftingRecipe struct {
	ID             uuid.UUID     `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	Name           string        `json:"name" gorm:"not null;index" validate:"required"`
	Description    string        `json:"description" gorm:"type:text"`
	RequiredItems  RequiredItems `json:"required_items" gorm:"type:json"`
	ResultItem     string        `json:"result_item" gorm:"not null"`
	ResultQuantity int           `json:"result_quantity" gorm:"default:1"`
	CraftingTime   int           `json:"crafting_time" gorm:"default:60"` // seconds
	RequiredLevel  int           `json:"required_level" gorm:"default:1"`
	IsActive       bool          `json:"is_active" gorm:"default:true"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`

	// Relationships
	CraftingSessions []CraftingSession `json:"crafting_sessions,omitempty" gorm:"foreignKey:RecipeID"`
}

func (cr *CraftingRecipe) BeforeCreate(tx *gorm.DB) error {
	if cr.ID == uuid.Nil {
		cr.ID = uuid.New()
	}
	return nil
}

type CraftingStatus string

const (
	CraftingStatusInProgress CraftingStatus = "in_progress"
	CraftingStatusCompleted  CraftingStatus = "completed"
	CraftingStatusCancelled  CraftingStatus = "cancelled"
)

type CraftingSession struct {
	ID          uuid.UUID      `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:char(36);not null;index"`
	RecipeID    uuid.UUID      `json:"recipe_id" gorm:"type:char(36);not null;index"`
	Status      CraftingStatus `json:"status" gorm:"type:enum('in_progress','completed','cancelled');default:'in_progress'"`
	StartedAt   time.Time      `json:"started_at"`
	CompletedAt *time.Time     `json:"completed_at"`

	// Relationships
	User   User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Recipe CraftingRecipe `json:"recipe,omitempty" gorm:"foreignKey:RecipeID"`
}

func (cs *CraftingSession) BeforeCreate(tx *gorm.DB) error {
	if cs.ID == uuid.Nil {
		cs.ID = uuid.New()
	}
	return nil
}
