package services

import (
	"errors"
	"time"

	"code-valley-api/internal/models"
	"code-valley-api/internal/repositories"
	"code-valley-api/internal/utils"

	"github.com/google/uuid"
)

type ShopService struct {
	shopRepo      *repositories.ShopRepository
	userRepo      *repositories.UserRepository
	inventoryRepo *repositories.InventoryRepository
}

func NewShopService() *ShopService {
	return &ShopService{
		shopRepo:      repositories.NewShopRepository(),
		userRepo:      repositories.NewUserRepository(),
		inventoryRepo: repositories.NewInventoryRepository(),
	}
}

func (s *ShopService) GetShopItems(pagination utils.PaginationParams) (*models.PaginatedResponse, error) {
	items, total, err := s.shopRepo.GetAllItems(pagination)
	if err != nil {
		return nil, err
	}

	data := make([]interface{}, len(items))
	for i, item := range items {
		data[i] = item
	}

	totalPages := int(total) / pagination.PerPage
	if int(total)%pagination.PerPage > 0 {
		totalPages++
	}

	return &models.PaginatedResponse{
		Data: data,
		Meta: models.PaginationMeta{
			CurrentPage: pagination.Page,
			PerPage:     pagination.PerPage,
			Total:       int(total),
			TotalPages:  totalPages,
		},
	}, nil
}

type BuyItemRequest struct {
	Quantity int `json:"quantity" validate:"min=1"`
}

func (s *ShopService) BuyItem(userID, itemID uuid.UUID, req BuyItemRequest) (*models.UserPurchase, error) {
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	// Get item
	item, err := s.shopRepo.GetItemByID(itemID)
	if err != nil {
		return nil, errors.New("item not found")
	}

	// Check stock
	if item.Stock != -1 && item.Stock < req.Quantity {
		return nil, errors.New("insufficient stock")
	}

	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	totalPrice := item.Price * req.Quantity

	// Check if user has enough coins
	if user.Coins < totalPrice {
		return nil, errors.New("insufficient coins")
	}

	// Create purchase
	purchase := &models.UserPurchase{
		UserID:     userID,
		ShopItemID: itemID,
		Quantity:   req.Quantity,
		TotalPrice: totalPrice,
		PurchasedAt: time.Now(),
	}

	if err := s.shopRepo.CreatePurchase(purchase); err != nil {
		return nil, err
	}

	// Update user coins
	user.Coins -= totalPrice
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Update item stock
	if item.Stock != -1 {
		newStock := item.Stock - req.Quantity
		if err := s.shopRepo.UpdateItemStock(itemID, newStock); err != nil {
			return nil, err
		}
	}

	// Add item to inventory
	inventory := &models.Inventory{
		UserID:   userID,
		ItemName: item.Name,
		Quantity: req.Quantity,
		ItemType: models.ItemType(item.ItemType),
	}
	s.inventoryRepo.AddItem(inventory)

	return purchase, nil
}

type SellItemRequest struct {
	Quantity int `json:"quantity" validate:"min=1"`
}

func (s *ShopService) SellItem(userID, itemID uuid.UUID, req SellItemRequest) (map[string]interface{}, error) {
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	// Get item from inventory
	inventoryItem, err := s.inventoryRepo.GetUserItem(userID, itemID)
	if err != nil {
		return nil, errors.New("item not found in inventory")
	}

	if inventoryItem.Quantity < req.Quantity {
		return nil, errors.New("insufficient quantity")
	}

	// Calculate sell price (50% of original price)
	sellPrice := 50 * req.Quantity // Base sell price

	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update user coins
	user.Coins += sellPrice
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Update inventory
	if inventoryItem.Quantity == req.Quantity {
		// Remove item completely
		if err := s.inventoryRepo.RemoveItem(userID, itemID); err != nil {
			return nil, err
		}
	} else {
		// Reduce quantity
		inventoryItem.Quantity -= req.Quantity
		if err := s.inventoryRepo.UpdateItem(inventoryItem); err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"coins_earned": sellPrice,
		"quantity_sold": req.Quantity,
	}, nil
}