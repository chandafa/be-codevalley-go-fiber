package handlers

import (
	"code-valley-api/internal/models"
	"code-valley-api/internal/services"
	"code-valley-api/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type FriendHandler struct {
	friendService *services.FriendService
}

func NewFriendHandler(friendService *services.FriendService) *FriendHandler {
	return &FriendHandler{
		friendService: friendService,
	}
}

func (h *FriendHandler) GetFriends(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)

	friends, err := h.friendService.GetUserFriends(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch friends"))
	}

	return c.JSON(models.SuccessResponse("Friends retrieved successfully", friends))
}

func (h *FriendHandler) SendFriendRequest(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)
	username := c.Params("username")

	err := h.friendService.SendFriendRequest(user.UserID, username)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("Friend request sent successfully", nil))
}

func (h *FriendHandler) AcceptFriendRequest(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)
	username := c.Params("username")

	err := h.friendService.AcceptFriendRequest(user.UserID, username)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("Friend request accepted successfully", nil))
}

func (h *FriendHandler) RemoveFriend(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)
	username := c.Params("username")

	err := h.friendService.RemoveFriend(user.UserID, username)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("Friend removed successfully", nil))
}

func (h *FriendHandler) GetOnlineFriends(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)

	friends, err := h.friendService.GetOnlineFriends(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch online friends"))
	}

	return c.JSON(models.SuccessResponse("Online friends retrieved successfully", friends))
}