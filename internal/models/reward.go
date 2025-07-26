package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DailyReward struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	Day         int       `json:"day" gorm:"not null;index"` // Day of the month or streak day
	RewardCoins int       `json:"reward_coins" gorm:"default:0"`
	RewardEXP   int       `json:"reward_exp" gorm:"default:0"`
	BonusItem   string    `json:"bonus_item"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (dr *DailyReward) BeforeCreate(tx *gorm.DB) error {
	if dr.ID == uuid.Nil {
		dr.ID = uuid.New()
	}
	return nil
}

type UserDailyReward struct {
	ID           uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index"`
	DailyRewardID uuid.UUID `json:"daily_reward_id" gorm:"type:char(36);not null;index"`
	ClaimedAt    time.Time `json:"claimed_at"`
	Date         time.Time `json:"date" gorm:"type:date;not null;index"`

	// Relationships
	User        User        `json:"user,omitempty" gorm:"foreignKey:UserID"`
	DailyReward DailyReward `json:"daily_reward,omitempty" gorm:"foreignKey:DailyRewardID"`
}

func (udr *UserDailyReward) BeforeCreate(tx *gorm.DB) error {
	if udr.ID == uuid.Nil {
		udr.ID = uuid.New()
	}
	return nil
}

type LoginStreak struct {
	ID           uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:char(36);not null;uniqueIndex"`
	CurrentStreak int      `json:"current_streak" gorm:"default:0"`
	LongestStreak int      `json:"longest_streak" gorm:"default:0"`
	LastLoginDate time.Time `json:"last_login_date" gorm:"type:date"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (ls *LoginStreak) BeforeCreate(tx *gorm.DB) error {
	if ls.ID == uuid.Nil {
		ls.ID = uuid.New()
	}
	return nil
}