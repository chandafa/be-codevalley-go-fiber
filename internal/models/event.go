package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventType string

const (
	EventTypeHackathon   EventType = "hackathon"
	EventTypeFestival    EventType = "festival"
	EventTypeCompetition EventType = "competition"
	EventTypeSeasonal    EventType = "seasonal"
)

type EventRewards map[string]interface{}

func (er EventRewards) Value() (driver.Value, error) {
	return json.Marshal(er)
}

func (er *EventRewards) Scan(value interface{}) error {
	if value == nil {
		*er = make(EventRewards)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, er)
}

type Event struct {
	ID           uuid.UUID    `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	Name         string       `json:"name" gorm:"not null;index" validate:"required"`
	Description  string       `json:"description" gorm:"type:text"`
	Type         EventType    `json:"type" gorm:"type:enum('hackathon','festival','competition','seasonal');not null"`
	StartDate    time.Time    `json:"start_date" gorm:"not null;index"`
	EndDate      time.Time    `json:"end_date" gorm:"not null;index"`
	Requirements EventRewards `json:"requirements" gorm:"type:json"`
	Rewards      EventRewards `json:"rewards" gorm:"type:json"`
	MaxParticipants int       `json:"max_participants" gorm:"default:-1"` // -1 means unlimited
	IsActive     bool         `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`

	// Relationships
	Participants []EventParticipant `json:"participants,omitempty" gorm:"foreignKey:EventID"`
}

func (e *Event) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

type ParticipantStatus string

const (
	ParticipantStatusJoined    ParticipantStatus = "joined"
	ParticipantStatusCompleted ParticipantStatus = "completed"
	ParticipantStatusDropped   ParticipantStatus = "dropped"
)

type EventParticipant struct {
	ID        uuid.UUID         `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID    uuid.UUID         `json:"user_id" gorm:"type:char(36);not null;index"`
	EventID   uuid.UUID         `json:"event_id" gorm:"type:char(36);not null;index"`
	Status    ParticipantStatus `json:"status" gorm:"type:enum('joined','completed','dropped');default:'joined'"`
	Score     int               `json:"score" gorm:"default:0"`
	JoinedAt  time.Time         `json:"joined_at"`
	CompletedAt *time.Time      `json:"completed_at"`

	// Relationships
	User  User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Event Event `json:"event,omitempty" gorm:"foreignKey:EventID"`
}

func (ep *EventParticipant) BeforeCreate(tx *gorm.DB) error {
	if ep.ID == uuid.Nil {
		ep.ID = uuid.New()
	}
	return nil
}