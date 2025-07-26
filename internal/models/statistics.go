package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserStatistics struct {
	ID                uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID            uuid.UUID `json:"user_id" gorm:"type:char(36);not null;uniqueIndex"`
	TotalPlayTime     int       `json:"total_play_time" gorm:"default:0"` // minutes
	QuestsCompleted   int       `json:"quests_completed" gorm:"default:0"`
	TasksCompleted    int       `json:"tasks_completed" gorm:"default:0"`
	CoinsEarned       int       `json:"coins_earned" gorm:"default:0"`
	CoinsSpent        int       `json:"coins_spent" gorm:"default:0"`
	FriendsCount      int       `json:"friends_count" gorm:"default:0"`
	GamesPlayed       int       `json:"games_played" gorm:"default:0"`
	GamesWon          int       `json:"games_won" gorm:"default:0"`
	ItemsCrafted      int       `json:"items_crafted" gorm:"default:0"`
	TutorialsCompleted int      `json:"tutorials_completed" gorm:"default:0"`
	LoginStreak       int       `json:"login_streak" gorm:"default:0"`
	LastActive        time.Time `json:"last_active"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (us *UserStatistics) BeforeCreate(tx *gorm.DB) error {
	if us.ID == uuid.Nil {
		us.ID = uuid.New()
	}
	return nil
}

type DailyStatistics struct {
	ID              uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID          uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index"`
	Date            time.Time `json:"date" gorm:"type:date;not null;index"`
	PlayTime        int       `json:"play_time" gorm:"default:0"` // minutes
	QuestsCompleted int       `json:"quests_completed" gorm:"default:0"`
	TasksCompleted  int       `json:"tasks_completed" gorm:"default:0"`
	CoinsEarned     int       `json:"coins_earned" gorm:"default:0"`
	EXPGained       int       `json:"exp_gained" gorm:"default:0"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (ds *DailyStatistics) BeforeCreate(tx *gorm.DB) error {
	if ds.ID == uuid.Nil {
		ds.ID = uuid.New()
	}
	return nil
}