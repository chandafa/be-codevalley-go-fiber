package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationType string

const (
	NotificationTypeQuest       NotificationType = "quest"
	NotificationTypeFriend      NotificationType = "friend"
	NotificationTypeAchievement NotificationType = "achievement"
	NotificationTypeSystem      NotificationType = "system"
	NotificationTypeNPC         NotificationType = "npc"
	NotificationTypeEvent       NotificationType = "event"
)

type Notification struct {
	ID       uuid.UUID        `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID   uuid.UUID        `json:"user_id" gorm:"type:char(36);not null;index"`
	Type     NotificationType `json:"type" gorm:"type:enum('quest','friend','achievement','system','npc','event');not null"`
	Title    string           `json:"title" gorm:"not null" validate:"required"`
	Message  string           `json:"message" gorm:"type:text" validate:"required"`
	IsRead   bool             `json:"is_read" gorm:"default:false"`
	Data     NotificationData `json:"data" gorm:"type:json"`
	CreatedAt time.Time       `json:"created_at"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type NotificationData map[string]interface{}

func (nd NotificationData) Value() (driver.Value, error) {
	return json.Marshal(nd)
}

func (nd *NotificationData) Scan(value interface{}) error {
	if value == nil {
		*nd = make(NotificationData)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, nd)
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}