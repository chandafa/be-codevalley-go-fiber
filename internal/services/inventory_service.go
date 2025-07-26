package services

import (
	"errors"

	"code-valley-api/internal/models"
	"code-valley-api/internal/repositories"

	"github.com/google/uuid"
)

type InventoryService struct {
	inventoryRepo *repositories.InventoryRepository
	userRepo      *repositories.UserRepository
}

func NewInventoryService() *InventoryService {
	return &InventoryService{
		inventoryRepo: repositories.NewInventoryRepository(),
		userRepo:      repositories.NewUserRepository(),
	}
}

func (s *InventoryService) GetUserInventory(userID uuid.UUID) ([]models.Inventory, error) {
	return s.inventoryRepo.GetUserInventory(userID)
}

func (s *InventoryService) GetUserItem(userID, itemID uuid.UUID) (*models.Inventory, error) {
	return s.inventoryRepo.GetUserItem(userID, itemID)
}

func (s *InventoryService) UseItem(userID, itemID uuid.UUID) (map[string]interface{}, error) {
	item, err := s.inventoryRepo.GetUserItem(userID, itemID)
	if err != nil {
		return nil, errors.New("item not found")
	}

	if item.Quantity <= 0 {
		return nil, errors.New("no items to use")
	}

	// Apply item effects based on item type
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	effects := s.applyItemEffects(user, item)

	// Reduce item quantity
	item.Quantity--
	if item.Quantity <= 0 {
		s.inventoryRepo.RemoveItem(userID, itemID)
	} else {
		s.inventoryRepo.UpdateItem(item)
	}

	// Update user
	s.userRepo.Update(user)

	return map[string]interface{}{
		"effects": effects,
		"remaining": item.Quantity,
	}, nil
}

func (s *InventoryService) EquipItem(userID, itemID uuid.UUID) (map[string]interface{}, error) {
	item, err := s.inventoryRepo.GetUserItem(userID, itemID)
	if err != nil {
		return nil, errors.New("item not found")
	}

	if item.ItemType != models.ItemTypeTool {
		return nil, errors.New("only tools can be equipped")
	}

	// In a real implementation, you'd have an equipped_items table
	// For now, we'll just return success
	return map[string]interface{}{
		"equipped": true,
		"item": item,
	}, nil
}

func (s *InventoryService) UnequipItem(userID, itemID uuid.UUID) (map[string]interface{}, error) {
	item, err := s.inventoryRepo.GetUserItem(userID, itemID)
	if err != nil {
		return nil, errors.New("item not found")
	}

	return map[string]interface{}{
		"unequipped": true,
		"item": item,
	}, nil
}

func (s *InventoryService) applyItemEffects(user *models.User, item *models.Inventory) map[string]interface{} {
	effects := make(map[string]interface{})

	switch item.ItemName {
	case "Coffee Beans":
		// Restore energy or give temporary boost
		effects["energy_boost"] = 10
	case "Health Potion":
		// Restore health
		effects["health_restored"] = 50
	case "EXP Boost":
		// Give EXP bonus
		bonus := 100
		user.EXP += bonus
		effects["exp_gained"] = bonus
	default:
		effects["message"] = "Item used successfully"
	}

	return effects
}