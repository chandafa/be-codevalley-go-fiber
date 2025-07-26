package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Guild struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	Name        string    `json:"name" gorm:"not null;uniqueIndex" validate:"required,min=3,max=50"`
	Description string    `json:"description" gorm:"type:text"`
	OwnerID     uuid.UUID `json:"owner_id" gorm:"type:char(36);not null;index"`
	MaxMembers  int       `json:"max_members" gorm:"default:20"`
	Level       int       `json:"level" gorm:"default:1"`
	EXP         int       `json:"exp" gorm:"default:0"`
	IconURL     string    `json:"icon_url"`
	IsPublic    bool      `json:"is_public" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Owner   User          `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	Members []GuildMember `json:"members,omitempty" gorm:"foreignKey:GuildID"`
}

func (g *Guild) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}

type GuildRole string

const (
	GuildRoleOwner     GuildRole = "owner"
	GuildRoleOfficer   GuildRole = "officer"
	GuildRoleMember    GuildRole = "member"
)

type GuildMember struct {
	ID       uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	GuildID  uuid.UUID `json:"guild_id" gorm:"type:char(36);not null;index"`
	UserID   uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index"`
	Role     GuildRole `json:"role" gorm:"type:enum('owner','officer','member');default:'member'"`
	JoinedAt time.Time `json:"joined_at"`

	// Relationships
	Guild Guild `json:"guild,omitempty" gorm:"foreignKey:GuildID"`
	User  User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (gm *GuildMember) BeforeCreate(tx *gorm.DB) error {
	if gm.ID == uuid.Nil {
		gm.ID = uuid.New()
	}
	return nil
}

type GuildInvitation struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	GuildID   uuid.UUID `json:"guild_id" gorm:"type:char(36);not null;index"`
	InviterID uuid.UUID `json:"inviter_id" gorm:"type:char(36);not null;index"`
	InviteeID uuid.UUID `json:"invitee_id" gorm:"type:char(36);not null;index"`
	Status    FriendshipStatus `json:"status" gorm:"type:enum('pending','accepted','blocked');default:'pending'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Guild   Guild `json:"guild,omitempty" gorm:"foreignKey:GuildID"`
	Inviter User  `json:"inviter,omitempty" gorm:"foreignKey:InviterID"`
	Invitee User  `json:"invitee,omitempty" gorm:"foreignKey:InviteeID"`
}

func (gi *GuildInvitation) BeforeCreate(tx *gorm.DB) error {
	if gi.ID == uuid.Nil {
		gi.ID = uuid.New()
	}
	return nil
}