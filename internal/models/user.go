package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole string

const (
	RolePlayer UserRole = "player"
	RoleAdmin  UserRole = "admin"
)

type User struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	Email       string    `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Username    string    `json:"username" gorm:"uniqueIndex;not null" validate:"required,min=3,max=30"`
	PasswordHash string   `json:"-" gorm:"not null"`
	Bio         string    `json:"bio" gorm:"type:text"`
	AvatarURL   string    `json:"avatar_url"`
	EXP         int       `json:"exp" gorm:"default:0"`
	Level       int       `json:"level" gorm:"default:1"`
	Coins       int       `json:"coins" gorm:"default:100"`
	Role        UserRole  `json:"role" gorm:"type:enum('player','admin');default:'player'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Inventory            []Inventory           `json:"inventory,omitempty" gorm:"foreignKey:UserID"`
	QuestProgress        []UserQuestProgress   `json:"quest_progress,omitempty" gorm:"foreignKey:UserID"`
	DailyTaskProgress    []UserDailyTaskProgress `json:"daily_task_progress,omitempty" gorm:"foreignKey:UserID"`
	Achievements         []UserAchievement     `json:"achievements,omitempty" gorm:"foreignKey:UserID"`
	StoryProgress        []StoryProgress       `json:"story_progress,omitempty" gorm:"foreignKey:UserID"`
	CodeBattles          []CodeBattle          `json:"code_battles,omitempty" gorm:"foreignKey:UserID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Bio       string    `json:"bio"`
	AvatarURL string    `json:"avatar_url"`
	EXP       int       `json:"exp"`
	Level     int       `json:"level"`
	Coins     int       `json:"coins"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		Bio:       u.Bio,
		AvatarURL: u.AvatarURL,
		EXP:       u.EXP,
		Level:     u.Level,
		Coins:     u.Coins,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}