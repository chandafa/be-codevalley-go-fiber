package main

import (
	"log"

	"code-valley-api/internal/config"
	"code-valley-api/internal/database"
	"code-valley-api/internal/middleware"
	"code-valley-api/internal/routes"
	"code-valley-api/internal/services"
	"code-valley-api/internal/websocket"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	if err := database.Initialize(cfg); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Run migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize WebSocket
	websocket.InitializeWebSocket()

	// Start game clock service
	gameClockService := services.NewGameClockService()
	gameClockService.Start()
	defer gameClockService.Stop()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"success": false,
				"message": err.Error(),
				"data":    nil,
			})
		},
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.CORSMiddleware(cfg))
	app.Use(middleware.RateLimitMiddleware(cfg))
	app.Use(middleware.ErrorHandlerMiddleware())

	// Setup routes
	routes.SetupRoutes(app, cfg)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}