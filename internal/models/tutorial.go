package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TutorialCategory string

const (
	TutorialCategoryBasics     TutorialCategory = "basics"
	TutorialCategoryAdvanced   TutorialCategory = "advanced"
	TutorialCategoryFramework  TutorialCategory = "framework"
	TutorialCategoryAlgorithm  TutorialCategory = "algorithm"
	TutorialCategoryBestPractice TutorialCategory = "best_practice"
)

type TutorialContent map[string]interface{}

func (tc TutorialContent) Value() (driver.Value, error) {
	return json.Marshal(tc)
}

func (tc *TutorialContent) Scan(value interface{}) error {
	if value == nil {
		*tc = make(TutorialContent)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, tc)
}

type Tutorial struct {
	ID           uuid.UUID        `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	Title        string           `json:"title" gorm:"not null;index" validate:"required"`
	Description  string           `json:"description" gorm:"type:text"`
	Category     TutorialCategory `json:"category" gorm:"type:enum('basics','advanced','framework','algorithm','best_practice');not null"`
	Content      TutorialContent  `json:"content" gorm:"type:json"`
	Prerequisites []uuid.UUID     `json:"prerequisites" gorm:"type:json"`
	RequiredLevel int             `json:"required_level" gorm:"default:1"`
	RewardEXP    int              `json:"reward_exp" gorm:"default:0"`
	EstimatedTime int             `json:"estimated_time" gorm:"default:30"` // minutes
	IsActive     bool             `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`

	// Relationships
	UserProgress []UserTutorialProgress `json:"user_progress,omitempty" gorm:"foreignKey:TutorialID"`
}

func (t *Tutorial) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

type TutorialStatus string

const (
	TutorialStatusNotStarted TutorialStatus = "not_started"
	TutorialStatusInProgress TutorialStatus = "in_progress"
	TutorialStatusCompleted  TutorialStatus = "completed"
)

type UserTutorialProgress struct {
	ID         uuid.UUID      `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID     uuid.UUID      `json:"user_id" gorm:"type:char(36);not null;index"`
	TutorialID uuid.UUID      `json:"tutorial_id" gorm:"type:char(36);not null;index"`
	Status     TutorialStatus `json:"status" gorm:"type:enum('not_started','in_progress','completed');default:'not_started'"`
	Progress   int            `json:"progress" gorm:"default:0"` // percentage
	StartedAt  *time.Time     `json:"started_at"`
	CompletedAt *time.Time    `json:"completed_at"`

	// Relationships
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Tutorial Tutorial `json:"tutorial,omitempty" gorm:"foreignKey:TutorialID"`
}

func (utp *UserTutorialProgress) BeforeCreate(tx *gorm.DB) error {
	if utp.ID == uuid.Nil {
		utp.ID = uuid.New()
	}
	return nil
}