package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NPCRelationship struct {
	ID             uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index"`
	NPCID          uuid.UUID `json:"npc_id" gorm:"type:char(36);not null;index"`
	FriendshipLevel int      `json:"friendship_level" gorm:"default:0"`
	LastInteraction time.Time `json:"last_interaction"`
	TotalInteractions int    `json:"total_interactions" gorm:"default:0"`
	GiftsGiven     int       `json:"gifts_given" gorm:"default:0"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	NPC  NPC  `json:"npc,omitempty" gorm:"foreignKey:NPCID"`
}

func (nr *NPCRelationship) BeforeCreate(tx *gorm.DB) error {
	if nr.ID == uuid.Nil {
		nr.ID = uuid.New()
	}
	return nil
}

type NPCInteraction struct {
	ID           uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index"`
	NPCID        uuid.UUID `json:"npc_id" gorm:"type:char(36);not null;index"`
	InteractionType string `json:"interaction_type" gorm:"not null"` // talk, gift, quest
	Data         InteractionData `json:"data" gorm:"type:json"`
	CreatedAt    time.Time `json:"created_at"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	NPC  NPC  `json:"npc,omitempty" gorm:"foreignKey:NPCID"`
}

type InteractionData map[string]interface{}

func (id InteractionData) Value() (driver.Value, error) {
	return json.Marshal(id)
}

func (id *InteractionData) Scan(value interface{}) error {
	if value == nil {
		*id = make(InteractionData)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, id)
}

func (ni *NPCInteraction) BeforeCreate(tx *gorm.DB) error {
	if ni.ID == uuid.Nil {
		ni.ID = uuid.New()
	}
	return nil
}