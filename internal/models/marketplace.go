package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MarketplaceItemStatus string

const (
	MarketplaceItemStatusActive MarketplaceItemStatus = "active"
	MarketplaceItemStatusSold   MarketplaceItemStatus = "sold"
	MarketplaceItemStatusExpired MarketplaceItemStatus = "expired"
)

type MarketplaceListing struct {
	ID          uuid.UUID             `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	SellerID    uuid.UUID             `json:"seller_id" gorm:"type:char(36);not null;index"`
	ItemName    string                `json:"item_name" gorm:"not null;index" validate:"required"`
	Description string                `json:"description" gorm:"type:text"`
	Price       int                   `json:"price" gorm:"not null" validate:"min=1"`
	Quantity    int                   `json:"quantity" gorm:"default:1" validate:"min=1"`
	ItemType    ItemType              `json:"item_type" gorm:"type:enum('tool','code','snippet','resource');not null"`
	Status      MarketplaceItemStatus `json:"status" gorm:"type:enum('active','sold','expired');default:'active'"`
	ExpiresAt   time.Time             `json:"expires_at"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`

	// Relationships
	Seller User `json:"seller,omitempty" gorm:"foreignKey:SellerID"`
}

func (ml *MarketplaceListing) BeforeCreate(tx *gorm.DB) error {
	if ml.ID == uuid.Nil {
		ml.ID = uuid.New()
	}
	return nil
}

type MarketplaceTransaction struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	ListingID uuid.UUID `json:"listing_id" gorm:"type:char(36);not null;index"`
	BuyerID   uuid.UUID `json:"buyer_id" gorm:"type:char(36);not null;index"`
	SellerID  uuid.UUID `json:"seller_id" gorm:"type:char(36);not null;index"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	TotalPrice int      `json:"total_price" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Listing MarketplaceListing `json:"listing,omitempty" gorm:"foreignKey:ListingID"`
	Buyer   User               `json:"buyer,omitempty" gorm:"foreignKey:BuyerID"`
	Seller  User               `json:"seller,omitempty" gorm:"foreignKey:SellerID"`
}

func (mt *MarketplaceTransaction) BeforeCreate(tx *gorm.DB) error {
	if mt.ID == uuid.Nil {
		mt.ID = uuid.New()
	}
	return nil
}