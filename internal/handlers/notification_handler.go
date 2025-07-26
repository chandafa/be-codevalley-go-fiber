package handlers

import (
	"code-valley-api/internal/models"
	"code-valley-api/internal/services"
	"code-valley-api/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)

	notifications, err := h.notificationService.GetUserNotifications(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to fetch notifications"))
	}

	return c.JSON(models.SuccessResponse("Notifications retrieved successfully", notifications))
}

func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)
	
	idParam := c.Params("id")
	notificationID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid notification ID"))
	}

	err = h.notificationService.MarkAsRead(user.UserID, notificationID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("Notification marked as read", nil))
}

func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)

	err := h.notificationService.MarkAllAsRead(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to mark notifications as read"))
	}

	return c.JSON(models.SuccessResponse("All notifications marked as read", nil))
}