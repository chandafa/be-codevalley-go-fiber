package handlers

import (
	"code-valley-api/internal/models"
	"code-valley-api/internal/services"
	"code-valley-api/internal/utils"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req services.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	response, err := h.authService.Register(req)
	if err != nil {
		if validationErrors := utils.FormatValidationErrors(err); len(validationErrors) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Message: "Validation failed",
				Data:    validationErrors,
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse("User registered successfully", response))
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req services.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	response, err := h.authService.Login(req)
	if err != nil {
		if validationErrors := utils.FormatValidationErrors(err); len(validationErrors) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Message: "Validation failed",
				Data:    validationErrors,
			})
		}
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse(err.Error()))
	}

	// return c.JSON(models.SuccessResponse("Login successful", response))
	return c.JSON(models.SuccessResponse("Login successful", map[string]interface{}{
		"token": response.Token,
		"user":  response.User,
	}))

}

func (h *AuthHandler) GetProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)

	profile, err := h.authService.GetProfile(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("User not found"))
	}

	return c.JSON(models.SuccessResponse("Profile retrieved successfully", profile))
}

func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)

	var req services.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid request body"))
	}

	profile, err := h.authService.UpdateProfile(user.UserID, req)
	if err != nil {
		if validationErrors := utils.FormatValidationErrors(err); len(validationErrors) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Message: "Validation failed",
				Data:    validationErrors,
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse(err.Error()))
	}

	return c.JSON(models.SuccessResponse("Profile updated successfully", profile))
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	// For simplicity, we'll just return the current user info
	// In a real application, you might want to implement proper refresh token logic
	user := c.Locals("user").(*utils.Claims)

	profile, err := h.authService.GetProfile(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.ErrorResponse("User not found"))
	}

	return c.JSON(models.SuccessResponse("Token refreshed successfully", map[string]interface{}{
		"user": profile,
	}))
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// In a real implementation, you might want to blacklist the token
	// For now, we'll just return success as the client should remove the token
	return c.JSON(models.SuccessResponse("Logged out successfully", nil))
}

func (h *AuthHandler) UploadAvatar(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)

	// Get uploaded file
	file, err := c.FormFile("avatar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("No file uploaded"))
	}

	// Validate file type
	if !isValidImageType(file) {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("Invalid file type. Only images are allowed"))
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse("File too large. Maximum size is 5MB"))
	}

	// In a real implementation, you would upload to cloud storage (AWS S3, etc.)
	// For now, we'll simulate by generating a URL
	avatarURL := "https://example.com/avatars/" + user.UserID.String() + ".jpg"

	// Update user profile
	profile, err := h.authService.UpdateAvatar(user.UserID, avatarURL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to update avatar"))
	}

	return c.JSON(models.SuccessResponse("Avatar updated successfully", profile))
}

func (h *AuthHandler) DeleteAccount(c *fiber.Ctx) error {
	user := c.Locals("user").(*utils.Claims)

	err := h.authService.DeleteAccount(user.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse("Failed to delete account"))
	}

	return c.JSON(models.SuccessResponse("Account deleted successfully", nil))
}

func isValidImageType(file *multipart.FileHeader) bool {
	validTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
	
	// Get content type from header
	contentType := file.Header.Get("Content-Type")
	return validTypes[contentType]
}
