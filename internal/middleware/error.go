package middleware

import (
	"log"

	"code-valley-api/internal/models"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandlerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			log.Printf("Error: %v", err)
			
			if fiberErr, ok := err.(*fiber.Error); ok {
				return c.Status(fiberErr.Code).JSON(models.ErrorResponse(fiberErr.Message))
			}
			
			return c.Status(fiber.StatusInternalServerError).JSON(
				models.ErrorResponse("Internal server error"),
			)
		}
		return nil
	}
}