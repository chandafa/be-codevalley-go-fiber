package handlers

import (
	"code-valley-api/internal/models"
	"code-valley-api/internal/services"
	"code-valley-api/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type InventoryHandler struct {
	inventoryService *services.InventoryService
}

func NewInventoryHandler(inventoryService *services.InventoryService) *InventoryHandler {
	return &InventoryHandler{
		inventoryService: inventoryService,
	}
}

func (h *InventoryHandler) GetInventory(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)

	inventory, err := h.inventoryService.GetUserInventory(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch inventory"))
	}

	return c.JSON(models.SuccessResponse("Inventory retrieved successfully", inventory))
}

func (h *InventoryHandler) GetItem(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)
	
	idParam := c.Params("id")
	itemID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid item ID"))
	}

	item, err := h.inventoryService.GetUserItem(user.UserID, itemID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("Item not found"))
	}

	return c.JSON(models.SuccessResponse("Item retrieved successfully", item))
}

func (h *InventoryHandler) UseItem(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)
	
	idParam := c.Params("id")
	itemID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid item ID"))
	}

	result, err := h.inventoryService.UseItem(user.UserID, itemID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("Item used successfully", result))
}

func (h *InventoryHandler) EquipItem(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)
	
	idParam := c.Params("id")
	itemID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid item ID"))
	}

	result, err := h.inventoryService.EquipItem(user.UserID, itemID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("Item equipped successfully", result))
}

func (h *InventoryHandler) UnequipItem(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)
	
	idParam := c.Params("id")
	itemID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid item ID"))
	}

	result, err := h.inventoryService.UnequipItem(user.UserID, itemID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("Item unequipped successfully", result))
}