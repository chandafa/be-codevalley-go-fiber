package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MiniGameType string

const (
	MiniGameTypeQuiz     MiniGameType = "quiz"
	MiniGameTypePuzzle   MiniGameType = "puzzle"
	MiniGameTypeRegex    MiniGameType = "regex"
	MiniGameTypeAlgorithm MiniGameType = "algorithm"
)

type GameConfig map[string]interface{}

func (gc GameConfig) Value() (driver.Value, error) {
	return json.Marshal(gc)
}

func (gc *GameConfig) Scan(value interface{}) error {
	if value == nil {
		*gc = make(GameConfig)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, gc)
}

type MiniGame struct {
	ID          uuid.UUID    `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	Name        string       `json:"name" gorm:"not null;index" validate:"required"`
	Description string       `json:"description" gorm:"type:text"`
	Type        MiniGameType `json:"type" gorm:"type:enum('quiz','puzzle','regex','algorithm');not null"`
	Difficulty  BattleDifficulty `json:"difficulty" gorm:"type:enum('easy','medium','hard');not null"`
	Config      GameConfig   `json:"config" gorm:"type:json"`
	RewardCoins int          `json:"reward_coins" gorm:"default:0"`
	RewardEXP   int          `json:"reward_exp" gorm:"default:0"`
	TimeLimit   int          `json:"time_limit" gorm:"default:300"` // seconds
	IsActive    bool         `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`

	// Relationships
	GameSessions []MiniGameSession `json:"game_sessions,omitempty" gorm:"foreignKey:MiniGameID"`
}

func (mg *MiniGame) BeforeCreate(tx *gorm.DB) error {
	if mg.ID == uuid.Nil {
		mg.ID = uuid.New()
	}
	return nil
}

type GameSessionStatus string

const (
	GameSessionStatusActive    GameSessionStatus = "active"
	GameSessionStatusCompleted GameSessionStatus = "completed"
	GameSessionStatusFailed    GameSessionStatus = "failed"
	GameSessionStatusTimeout   GameSessionStatus = "timeout"
)

type MiniGameSession struct {
	ID         uuid.UUID         `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID     uuid.UUID         `json:"user_id" gorm:"type:char(36);not null;index"`
	MiniGameID uuid.UUID         `json:"mini_game_id" gorm:"type:char(36);not null;index"`
	Status     GameSessionStatus `json:"status" gorm:"type:enum('active','completed','failed','timeout');default:'active'"`
	Score      int               `json:"score" gorm:"default:0"`
	StartedAt  time.Time         `json:"started_at"`
	CompletedAt *time.Time       `json:"completed_at"`
	SessionData GameConfig       `json:"session_data" gorm:"type:json"`

	// Relationships
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	MiniGame MiniGame `json:"mini_game,omitempty" gorm:"foreignKey:MiniGameID"`
}

func (mgs *MiniGameSession) BeforeCreate(tx *gorm.DB) error {
	if mgs.ID == uuid.Nil {
		mgs.ID = uuid.New()
	}
	return nil
}