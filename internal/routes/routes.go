package routes

import (
	"code-valley-api/internal/config"
	"code-valley-api/internal/handlers"
	"code-valley-api/internal/middleware"
	"code-valley-api/internal/services"
	"code-valley-api/internal/websocket"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, cfg *config.Config) {
	// Initialize services
	authService := services.NewAuthService(cfg)
	questService := services.NewQuestService()
	friendService := services.NewFriendService()
	leaderboardService := services.NewLeaderboardService()
	shopService := services.NewShopService()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	questHandler := handlers.NewQuestHandler(questService)
	friendHandler := handlers.NewFriendHandler(friendService)
	leaderboardHandler := handlers.NewLeaderboardHandler(leaderboardService)
	shopHandler := handlers.NewShopHandler(shopService)

	// WebSocket endpoint
	app.Get("/ws", websocket.WebSocketUpgrade(cfg))

	// API v1 group
	api := app.Group("/api/v1")

	// Public auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected auth routes
	authProtected := auth.Group("/", middleware.AuthMiddleware(cfg))
	authProtected.Get("/me", authHandler.GetProfile)
	authProtected.Put("/profile", authHandler.UpdateProfile)
	authProtected.Post("/refresh", authHandler.RefreshToken)

	// Protected quest routes
	quests := api.Group("/quests", middleware.AuthMiddleware(cfg))
	quests.Get("/", questHandler.GetQuests)
	quests.Get("/:id", questHandler.GetQuestByID)
	quests.Post("/:id/start", questHandler.StartQuest)
	quests.Post("/:id/complete", questHandler.CompleteQuest)
	quests.Get("/progress", questHandler.GetUserProgress)

	// Admin quest routes
	adminQuests := api.Group("/admin/quests", middleware.AuthMiddleware(cfg), middleware.RequireRole("admin"))
	adminQuests.Post("/", questHandler.CreateQuest)
	adminQuests.Put("/:id", questHandler.UpdateQuest)
	adminQuests.Delete("/:id", questHandler.DeleteQuest)

	// Friend routes
	friends := api.Group("/friends", middleware.AuthMiddleware(cfg))
	friends.Get("/", friendHandler.GetFriends)
	friends.Post("/:username/add", friendHandler.SendFriendRequest)
	friends.Post("/:username/accept", friendHandler.AcceptFriendRequest)
	friends.Delete("/:username/remove", friendHandler.RemoveFriend)
	friends.Get("/online", friendHandler.GetOnlineFriends)

	// Leaderboard routes
	leaderboard := api.Group("/leaderboard")
	leaderboard.Get("/coins", leaderboardHandler.GetCoinLeaderboard)
	leaderboard.Get("/exp", leaderboardHandler.GetEXPLeaderboard)
	leaderboard.Get("/tasks", leaderboardHandler.GetTaskLeaderboard)

	// Shop routes
	shop := api.Group("/shop", middleware.AuthMiddleware(cfg))
	shop.Get("/items", shopHandler.GetShopItems)
	shop.Post("/items/:id/buy", shopHandler.BuyItem)
	shop.Post("/items/:id/sell", shopHandler.SellItem)

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "Code Valley API is running",
		})
	})
}