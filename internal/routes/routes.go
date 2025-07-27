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
	notificationService := services.NewNotificationService()
	inventoryService := services.NewInventoryService()
	adminService := services.NewAdminService()
	worldService := services.NewWorldService()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	questHandler := handlers.NewQuestHandler(questService)
	friendHandler := handlers.NewFriendHandler(friendService)
	leaderboardHandler := handlers.NewLeaderboardHandler(leaderboardService)
	shopHandler := handlers.NewShopHandler(shopService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)
	adminHandler := handlers.NewAdminHandler(adminService)
	worldHandler := handlers.NewWorldHandler(worldService)

	// WebSocket endpoint
	app.Get("/ws", websocket.WebSocketUpgrade(cfg))

	// API v1 group
	api := app.Group("/api/v1")

	// World/Map routes
	world := api.Group("/world", middleware.AuthMiddleware(cfg))
	world.Get("/maps/:map_name/state", worldHandler.GetMapState)
	world.Get("/position", worldHandler.GetPlayerPosition)
	world.Post("/teleport", worldHandler.TeleportPlayer)
	world.Get("/time", worldHandler.GetGameTime)
	
	// Code farming routes
	farming := api.Group("/farming", middleware.AuthMiddleware(cfg))
	farming.Get("/", worldHandler.GetCodeFarms)
	farming.Post("/plant", worldHandler.PlantCode)
	farming.Post("/:id/water", worldHandler.WaterCode)
	farming.Post("/:id/harvest", worldHandler.HarvestCode)

	// Public auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// Protected auth routes
	authProtected := auth.Group("/", middleware.AuthMiddleware(cfg))
	authProtected.Get("/me", authHandler.GetProfile)
	authProtected.Put("/profile", authHandler.UpdateProfile)
	authProtected.Post("/refresh", authHandler.RefreshToken)
	authProtected.Post("/logout", authHandler.Logout)
	authProtected.Post("/avatar", authHandler.UploadAvatar)
	authProtected.Delete("/delete", authHandler.DeleteAccount)

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

	// Notification routes
	notifications := api.Group("/notifications", middleware.AuthMiddleware(cfg))
	notifications.Get("/", notificationHandler.GetNotifications)
	notifications.Post("/:id/read", notificationHandler.MarkAsRead)
	notifications.Post("/mark-read", notificationHandler.MarkAllAsRead)

	// Inventory routes
	inventory := api.Group("/inventory", middleware.AuthMiddleware(cfg))
	inventory.Get("/", inventoryHandler.GetInventory)
	inventory.Get("/:id", inventoryHandler.GetItem)
	inventory.Post("/use/:id", inventoryHandler.UseItem)
	inventory.Post("/equip/:id", inventoryHandler.EquipItem)
	inventory.Post("/unequip/:id", inventoryHandler.UnequipItem)

	// Admin routes
	admin := api.Group("/admin", middleware.AuthMiddleware(cfg), middleware.RequireRole("admin"))
	admin.Get("/users", adminHandler.GetAllUsers)
	admin.Put("/users/:id/ban", adminHandler.BanUser)
	admin.Put("/users/:id/role", adminHandler.ChangeUserRole)
	admin.Get("/stats", adminHandler.GetSystemStats)
	admin.Get("/logs", adminHandler.GetAuditLogs)

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "Code Valley API is running",
		})
	})
}