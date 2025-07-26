package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemType string

const (
	ItemTypeTool     ItemType = "tool"
	ItemTypeCode     ItemType = "code"
	ItemTypeSnippet  ItemType = "snippet"
	ItemTypeResource ItemType = "resource"
)

type Inventory struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index"`
	ItemName  string    `json:"item_name" gorm:"not null" validate:"required"`
	Quantity  int       `json:"quantity" gorm:"default:1" validate:"min=1"`
	ItemType  ItemType  `json:"item_type" gorm:"type:enum('tool','code','snippet','resource');not null" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (i *Inventory) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}