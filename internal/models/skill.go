package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SkillCategory string

const (
	SkillCategoryProgramming SkillCategory = "programming"
	SkillCategoryDebugging   SkillCategory = "debugging"
	SkillCategoryOptimization SkillCategory = "optimization"
	SkillCategorySocial      SkillCategory = "social"
	SkillCategoryEconomic    SkillCategory = "economic"
)

type SkillEffects map[string]interface{}

func (se SkillEffects) Value() (driver.Value, error) {
	return json.Marshal(se)
}

func (se *SkillEffects) Scan(value interface{}) error {
	if value == nil {
		*se = make(SkillEffects)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, se)
}

type Skill struct {
	ID           uuid.UUID     `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	Name         string        `json:"name" gorm:"not null;index" validate:"required"`
	Description  string        `json:"description" gorm:"type:text"`
	Category     SkillCategory `json:"category" gorm:"type:enum('programming','debugging','optimization','social','economic');not null"`
	MaxLevel     int           `json:"max_level" gorm:"default:10"`
	BaseCost     int           `json:"base_cost" gorm:"default:100"`
	Effects      SkillEffects  `json:"effects" gorm:"type:json"`
	IconURL      string        `json:"icon_url"`
	Prerequisites []uuid.UUID  `json:"prerequisites" gorm:"type:json"`
	IsActive     bool          `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`

	// Relationships
	UserSkills []UserSkill `json:"user_skills,omitempty" gorm:"foreignKey:SkillID"`
}

func (s *Skill) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type UserSkill struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index"`
	SkillID   uuid.UUID `json:"skill_id" gorm:"type:char(36);not null;index"`
	Level     int       `json:"level" gorm:"default:1"`
	UnlockedAt time.Time `json:"unlocked_at"`

	// Relationships
	User  User  `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Skill Skill `json:"skill,omitempty" gorm:"foreignKey:SkillID"`
}

func (us *UserSkill) BeforeCreate(tx *gorm.DB) error {
	if us.ID == uuid.Nil {
		us.ID = uuid.New()
	}
	return nil
}