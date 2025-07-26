package handlers

import (
	"code-valley-api/internal/models"
	"code-valley-api/internal/services"
	"code-valley-api/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ShopHandler struct {
	shopService *services.ShopService
}

func NewShopHandler(shopService *services.ShopService) *ShopHandler {
	return &ShopHandler{
		shopService: shopService,
	}
}

func (h *ShopHandler) GetShopItems(c *fiber.Ctx) error {
	pagination := utils.GetPaginationParams(c)

	response, err := h.shopService.GetShopItems(pagination)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch shop items"))
	}

	return c.JSON(models.SuccessResponse("Shop items retrieved successfully", response))
}

func (h *ShopHandler) BuyItem(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)
	
	idParam := c.Params("id")
	itemID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid item ID"))
	}

	var req services.BuyItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	purchase, err := h.shopService.BuyItem(user.UserID, itemID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("Item purchased successfully", purchase))
}

func (h *ShopHandler) SellItem(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)
	
	idParam := c.Params("id")
	itemID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid item ID"))
	}

	var req services.SellItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	result, err := h.shopService.SellItem(user.UserID, itemID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("Item sold successfully", result))
}