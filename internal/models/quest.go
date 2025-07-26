package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RequiredItems map[string]int

func (ri RequiredItems) Value() (driver.Value, error) {
	return json.Marshal(ri)
}

func (ri *RequiredItems) Scan(value interface{}) error {
	if value == nil {
		*ri = make(RequiredItems)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, ri)
}

type Quest struct {
	ID            uuid.UUID     `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	Title         string        `json:"title" gorm:"not null" validate:"required"`
	Description   string        `json:"description" gorm:"type:text" validate:"required"`
	RewardCoins   int           `json:"reward_coins" gorm:"default:0"`
	RewardEXP     int           `json:"reward_exp" gorm:"default:0"`
	RequiredItems RequiredItems `json:"required_items" gorm:"type:json"`
	IsRepeatable  bool          `json:"is_repeatable" gorm:"default:false"`
	IsActive      bool          `json:"is_active" gorm:"default:true"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`

	// Relationships
	UserProgress []UserQuestProgress `json:"user_progress,omitempty" gorm:"foreignKey:QuestID"`
}

func (q *Quest) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}

type QuestStatus string

const (
	QuestStatusNotStarted QuestStatus = "not_started"
	QuestStatusInProgress QuestStatus = "in_progress"
	QuestStatusCompleted  QuestStatus = "completed"
)

type ProgressData map[string]interface{}

func (pd ProgressData) Value() (driver.Value, error) {
	return json.Marshal(pd)
}

func (pd *ProgressData) Scan(value interface{}) error {
	if value == nil {
		*pd = make(ProgressData)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, pd)
}

type UserQuestProgress struct {
	ID           uuid.UUID     `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID       uuid.UUID     `json:"user_id" gorm:"type:char(36);not null;index"`
	QuestID      uuid.UUID     `json:"quest_id" gorm:"type:char(36);not null;index"`
	Status       QuestStatus   `json:"status" gorm:"type:enum('not_started','in_progress','completed');default:'not_started'"`
	ProgressData ProgressData  `json:"progress_data" gorm:"type:json"`
	StartedAt    *time.Time    `json:"started_at"`
	CompletedAt  *time.Time    `json:"completed_at"`

	// Relationships
	User  User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Quest Quest `json:"quest,omitempty" gorm:"foreignKey:QuestID"`
}

func (uqp *UserQuestProgress) BeforeCreate(tx *gorm.DB) error {
	if uqp.ID == uuid.Nil {
		uqp.ID = uuid.New()
	}
	return nil
}