package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShopItemType string

const (
	ShopItemTypeTool      ShopItemType = "tool"
	ShopItemTypeUpgrade   ShopItemType = "upgrade"
	ShopItemTypeCosmetic  ShopItemType = "cosmetic"
	ShopItemTypeResource  ShopItemType = "resource"
	ShopItemTypeSkill     ShopItemType = "skill"
)

type ShopItem struct {
	ID          uuid.UUID    `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	Name        string       `json:"name" gorm:"not null;index" validate:"required"`
	Description string       `json:"description" gorm:"type:text"`
	Price       int          `json:"price" gorm:"not null" validate:"min=0"`
	ItemType    ShopItemType `json:"item_type" gorm:"type:enum('tool','upgrade','cosmetic','resource','skill');not null"`
	IconURL     string       `json:"icon_url"`
	IsAvailable bool         `json:"is_available" gorm:"default:true"`
	Stock       int          `json:"stock" gorm:"default:-1"` // -1 means unlimited
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (si *ShopItem) BeforeCreate(tx *gorm.DB) error {
	if si.ID == uuid.Nil {
		si.ID = uuid.New()
	}
	return nil
}

type UserPurchase struct {
	ID         uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:char(36);not null;index"`
	ShopItemID uuid.UUID `json:"shop_item_id" gorm:"type:char(36);not null;index"`
	Quantity   int       `json:"quantity" gorm:"default:1"`
	TotalPrice int       `json:"total_price" gorm:"not null"`
	PurchasedAt time.Time `json:"purchased_at"`

	// Relationships
	User     User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	ShopItem ShopItem `json:"shop_item,omitempty" gorm:"foreignKey:ShopItemID"`
}

func (up *UserPurchase) BeforeCreate(tx *gorm.DB) error {
	if up.ID == uuid.Nil {
		up.ID = uuid.New()
	}
	return nil
}