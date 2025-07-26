package repositories

import (
	"code-valley-api/internal/database"
	"code-valley-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository() *InventoryRepository {
	return &InventoryRepository{
		db: database.GetDB(),
	}
}

func (r *InventoryRepository) AddItem(item *models.Inventory) error {
	// Check if item already exists
	var existing models.Inventory
	err := r.db.Where("user_id = ? AND item_name = ?", item.UserID, item.ItemName).First(&existing).Error
	
	if err == nil {
		// Item exists, update quantity
		existing.Quantity += item.Quantity
		return r.db.Save(&existing).Error
	} else if err == gorm.ErrRecordNotFound {
		// Item doesn't exist, create new
		return r.db.Create(item).Error
	}
	
	return err
}

func (r *InventoryRepository) GetUserItem(userID, itemID uuid.UUID) (*models.Inventory, error) {
	var item models.Inventory
	err := r.db.Where("user_id = ? AND id = ?", userID, itemID).First(&item).Error
	return &item, err
}

func (r *InventoryRepository) UpdateItem(item *models.Inventory) error {
	return r.db.Save(item).Error
}

func (r *InventoryRepository) RemoveItem(userID, itemID uuid.UUID) error {
	return r.db.Where("user_id = ? AND id = ?", userID, itemID).Delete(&models.Inventory{}).Error
}

func (r *InventoryRepository) GetUserInventory(userID uuid.UUID) ([]models.Inventory, error) {
	var items []models.Inventory
	err := r.db.Where("user_id = ?", userID).Find(&items).Error
	return items, err
}