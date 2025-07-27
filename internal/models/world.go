package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MapType string

const (
	MapTypeVillage    MapType = "village"
	MapTypeCodeMine   MapType = "code_mine"
	MapTypeDataFarm   MapType = "data_farm"
	MapTypeOffice     MapType = "office"
	MapTypeLibrary    MapType = "library"
	MapTypeLab        MapType = "lab"
)

type Map struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	Name        string    `json:"name" gorm:"not null;uniqueIndex" validate:"required"`
	Type        MapType   `json:"type" gorm:"type:enum('village','code_mine','data_farm','office','library','lab');not null"`
	Width       int       `json:"width" gorm:"not null" validate:"min=1"`
	Height      int       `json:"height" gorm:"not null" validate:"min=1"`
	Layout      MapLayout `json:"layout" gorm:"type:json"`
	Description string    `json:"description" gorm:"type:text"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	PlayerPositions []PlayerPosition `json:"player_positions,omitempty" gorm:"foreignKey:MapID"`
	WorldObjects    []WorldObject    `json:"world_objects,omitempty" gorm:"foreignKey:MapID"`
	NPCPositions    []NPCPosition    `json:"npc_positions,omitempty" gorm:"foreignKey:MapID"`
}

type MapLayout map[string]interface{}

func (ml MapLayout) Value() (driver.Value, error) {
	return json.Marshal(ml)
}

func (ml *MapLayout) Scan(value interface{}) error {
	if value == nil {
		*ml = make(MapLayout)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, ml)
}

func (m *Map) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

type PlayerPosition struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:char(36);not null;uniqueIndex"`
	MapID     uuid.UUID `json:"map_id" gorm:"type:char(36);not null;index"`
	PosX      int       `json:"pos_x" gorm:"not null"`
	PosY      int       `json:"pos_y" gorm:"not null"`
	Direction string    `json:"direction" gorm:"default:'down'"` // up, down, left, right
	LastMoved time.Time `json:"last_moved"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Map  Map  `json:"map,omitempty" gorm:"foreignKey:MapID"`
}

func (pp *PlayerPosition) BeforeCreate(tx *gorm.DB) error {
	if pp.ID == uuid.Nil {
		pp.ID = uuid.New()
	}
	return nil
}

type ObjectType string

const (
	ObjectTypeTree      ObjectType = "tree"
	ObjectTypeRock      ObjectType = "rock"
	ObjectTypeChest     ObjectType = "chest"
	ObjectTypeServer    ObjectType = "server"
	ObjectTypeWorkstation ObjectType = "workstation"
	ObjectTypeCodeBlock ObjectType = "code_block"
	ObjectTypeBugHive   ObjectType = "bug_hive"
)

type ObjectState map[string]interface{}

func (os ObjectState) Value() (driver.Value, error) {
	return json.Marshal(os)
}

func (os *ObjectState) Scan(value interface{}) error {
	if value == nil {
		*os = make(ObjectState)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, os)
}

type WorldObject struct {
	ID         uuid.UUID   `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	MapID      uuid.UUID   `json:"map_id" gorm:"type:char(36);not null;index"`
	ObjectType ObjectType  `json:"object_type" gorm:"type:enum('tree','rock','chest','server','workstation','code_block','bug_hive');not null"`
	PosX       int         `json:"pos_x" gorm:"not null"`
	PosY       int         `json:"pos_y" gorm:"not null"`
	State      ObjectState `json:"state" gorm:"type:json"`
	IsActive   bool        `json:"is_active" gorm:"default:true"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`

	// Relationships
	Map Map `json:"map,omitempty" gorm:"foreignKey:MapID"`
}

func (wo *WorldObject) BeforeCreate(tx *gorm.DB) error {
	if wo.ID == uuid.Nil {
		wo.ID = uuid.New()
	}
	return nil
}

type NPCPosition struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	NPCID     uuid.UUID `json:"npc_id" gorm:"type:char(36);not null;index"`
	MapID     uuid.UUID `json:"map_id" gorm:"type:char(36);not null;index"`
	PosX      int       `json:"pos_x" gorm:"not null"`
	PosY      int       `json:"pos_y" gorm:"not null"`
	Direction string    `json:"direction" gorm:"default:'down'"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	NPC NPC `json:"npc,omitempty" gorm:"foreignKey:NPCID"`
	Map Map `json:"map,omitempty" gorm:"foreignKey:MapID"`
}

func (np *NPCPosition) BeforeCreate(tx *gorm.DB) error {
	if np.ID == uuid.Nil {
		np.ID = uuid.New()
	}
	return nil
}

type NPCSchedule struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	NPCID     uuid.UUID `json:"npc_id" gorm:"type:char(36);not null;index"`
	DayOfWeek int       `json:"day_of_week" gorm:"not null"` // 0-6 (Sunday-Saturday)
	TimeOfDay int       `json:"time_of_day" gorm:"not null"` // 0-2359 (24-hour format)
	MapID     uuid.UUID `json:"map_id" gorm:"type:char(36);not null"`
	PosX      int       `json:"pos_x" gorm:"not null"`
	PosY      int       `json:"pos_y" gorm:"not null"`
	Action    string    `json:"action"` // "work", "rest", "patrol", etc.
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	NPC NPC `json:"npc,omitempty" gorm:"foreignKey:NPCID"`
	Map Map `json:"map,omitempty" gorm:"foreignKey:MapID"`
}

func (ns *NPCSchedule) BeforeCreate(tx *gorm.DB) error {
	if ns.ID == uuid.Nil {
		ns.ID = uuid.New()
	}
	return nil
}

type GameClock struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	GameYear    int       `json:"game_year" gorm:"default:1"`
	GameSeason  string    `json:"game_season" gorm:"default:'spring'"` // spring, summer, fall, winter
	GameDay     int       `json:"game_day" gorm:"default:1"`
	GameHour    int       `json:"game_hour" gorm:"default:6"`
	GameMinute  int       `json:"game_minute" gorm:"default:0"`
	IsPaused    bool      `json:"is_paused" gorm:"default:false"`
	TimeScale   float64   `json:"time_scale" gorm:"default:1.0"` // 1.0 = normal speed
	UpdatedAt   time.Time `json:"updated_at"`
}

func (gc *GameClock) BeforeCreate(tx *gorm.DB) error {
	if gc.ID == uuid.Nil {
		gc.ID = uuid.New()
	}
	return nil
}

type CodeFarm struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index"`
	PlotX       int       `json:"plot_x" gorm:"not null"`
	PlotY       int       `json:"plot_y" gorm:"not null"`
	CodeType    string    `json:"code_type"` // "algorithm", "function", "class", etc.
	PlantedAt   *time.Time `json:"planted_at"`
	LastWatered *time.Time `json:"last_watered"`
	HarvestAt   *time.Time `json:"harvest_at"`
	GrowthStage int       `json:"growth_stage" gorm:"default:0"` // 0-4
	Quality     string    `json:"quality" gorm:"default:'normal'"` // normal, silver, gold, iridium
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (cf *CodeFarm) BeforeCreate(tx *gorm.DB) error {
	if cf.ID == uuid.Nil {
		cf.ID = uuid.New()
	}
	return nil
}