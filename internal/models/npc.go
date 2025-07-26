package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NPCRole string

const (
	NPCRoleMentor   NPCRole = "mentor"
	NPCRoleClient   NPCRole = "client"
	NPCRoleVillager NPCRole = "villager"
)

type QuestIDs []uuid.UUID

func (qi QuestIDs) Value() (driver.Value, error) {
	return json.Marshal(qi)
}

func (qi *QuestIDs) Scan(value interface{}) error {
	if value == nil {
		*qi = []uuid.UUID{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, qi)
}

type NPC struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	Name        string    `json:"name" gorm:"not null" validate:"required"`
	Role        NPCRole   `json:"role" gorm:"type:enum('mentor','client','villager');not null" validate:"required"`
	Dialogue    string    `json:"dialogue" gorm:"type:text"`
	Location    string    `json:"location" gorm:"not null" validate:"required"`
	AvatarURL   string    `json:"avatar_url"`
	QuestsGiven QuestIDs  `json:"quests_given" gorm:"type:json"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (n *NPC) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}