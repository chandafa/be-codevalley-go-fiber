package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BattleDifficulty string

const (
	DifficultyEasy   BattleDifficulty = "easy"
	DifficultyMedium BattleDifficulty = "medium"
	DifficultyHard   BattleDifficulty = "hard"
)

type BattleStatus string

const (
	BattleStatusInProgress BattleStatus = "in_progress"
	BattleStatusCompleted  BattleStatus = "completed"
	BattleStatusFailed     BattleStatus = "failed"
)

type CodeBattle struct {
	ID            uuid.UUID        `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID        uuid.UUID        `json:"user_id" gorm:"type:char(36);not null;index"`
	ChallengeName string           `json:"challenge_name" gorm:"not null" validate:"required"`
	Difficulty    BattleDifficulty `json:"difficulty" gorm:"type:enum('easy','medium','hard');not null" validate:"required"`
	Status        BattleStatus     `json:"status" gorm:"type:enum('in_progress','completed','failed');default:'in_progress'"`
	Score         int              `json:"score" gorm:"default:0"`
	StartedAt     time.Time        `json:"started_at"`
	CompletedAt   *time.Time       `json:"completed_at"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (cb *CodeBattle) BeforeCreate(tx *gorm.DB) error {
	if cb.ID == uuid.Nil {
		cb.ID = uuid.New()
	}
	return nil
}