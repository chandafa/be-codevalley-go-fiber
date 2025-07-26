package repositories

import (
	"code-valley-api/internal/database"
	"code-valley-api/internal/models"
	"code-valley-api/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShopRepository struct {
	db *gorm.DB
}

func NewShopRepository() *ShopRepository {
	return &ShopRepository{
		db: database.GetDB(),
	}
}

func (r *ShopRepository) GetAllItems(pagination utils.PaginationParams) ([]models.ShopItem, int64, error) {
	var items []models.ShopItem
	var total int64

	query := r.db.Model(&models.ShopItem{}).Where("is_available = ?", true)
	query.Count(&total)

	err := query.Offset(pagination.Offset).
		Limit(pagination.PerPage).
		Find(&items).Error

	return items, total, err
}

func (r *ShopRepository) GetItemByID(id uuid.UUID) (*models.ShopItem, error) {
	var item models.ShopItem
	err := r.db.First(&item, "id = ? AND is_available = ?", id, true).Error
	return &item, err
}

func (r *ShopRepository) CreatePurchase(purchase *models.UserPurchase) error {
	return r.db.Create(purchase).Error
}

func (r *ShopRepository) UpdateItemStock(itemID uuid.UUID, newStock int) error {
	return r.db.Model(&models.ShopItem{}).Where("id = ?", itemID).Update("stock", newStock).Error
}

func (r *ShopRepository) GetUserPurchases(userID uuid.UUID) ([]models.UserPurchase, error) {
	var purchases []models.UserPurchase
	err := r.db.Preload("ShopItem").Where("user_id = ?", userID).Find(&purchases).Error
	return purchases, err
}