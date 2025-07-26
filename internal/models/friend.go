package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FriendshipStatus string

const (
	FriendshipStatusPending  FriendshipStatus = "pending"
	FriendshipStatusAccepted FriendshipStatus = "accepted"
	FriendshipStatusBlocked  FriendshipStatus = "blocked"
)

type Friendship struct {
	ID           uuid.UUID        `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	RequesterID  uuid.UUID        `json:"requester_id" gorm:"type:char(36);not null;index"`
	AddresseeID  uuid.UUID        `json:"addressee_id" gorm:"type:char(36);not null;index"`
	Status       FriendshipStatus `json:"status" gorm:"type:enum('pending','accepted','blocked');default:'pending'"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`

	// Relationships
	Requester User `json:"requester,omitempty" gorm:"foreignKey:RequesterID"`
	Addressee User `json:"addressee,omitempty" gorm:"foreignKey:AddresseeID"`
}

func (f *Friendship) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

type OnlineUser struct {
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:char(36);not null;uniqueIndex"`
	LastSeen   time.Time `json:"last_seen"`
	IsOnline   bool      `json:"is_online" gorm:"default:true"`
	SocketID   string    `json:"socket_id" gorm:"index"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (ou *OnlineUser) BeforeCreate(tx *gorm.DB) error {
	if ou.ID == uuid.Nil {
		ou.ID = uuid.New()
	}
	return nil
}