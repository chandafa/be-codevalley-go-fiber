package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DailyTask struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	TaskName    string    `json:"task_name" gorm:"not null" validate:"required"`
	TaskType    string    `json:"task_type" gorm:"not null" validate:"required"`
	Description string    `json:"description" gorm:"type:text" validate:"required"`
	RewardEXP   int       `json:"reward_exp" gorm:"default:0"`
	RewardCoins int       `json:"reward_coins" gorm:"default:0"`
	Date        time.Time `json:"date" gorm:"type:date;not null"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	UserProgress []UserDailyTaskProgress `json:"user_progress,omitempty" gorm:"foreignKey:DailyTaskID"`
}

func (dt *DailyTask) BeforeCreate(tx *gorm.DB) error {
	if dt.ID == uuid.Nil {
		dt.ID = uuid.New()
	}
	return nil
}

type UserDailyTaskProgress struct {
	ID          uuid.UUID  `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID      uuid.UUID  `json:"user_id" gorm:"type:char(36);not null;index"`
	DailyTaskID uuid.UUID  `json:"daily_task_id" gorm:"type:char(36);not null;index"`
	CompletedAt *time.Time `json:"completed_at"`
	Date        time.Time  `json:"date" gorm:"type:date;not null"`

	// Relationships
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	DailyTask DailyTask `json:"daily_task,omitempty" gorm:"foreignKey:DailyTaskID"`
}

func (udtp *UserDailyTaskProgress) BeforeCreate(tx *gorm.DB) error {
	if udtp.ID == uuid.Nil {
		udtp.ID = uuid.New()
	}
	return nil
}