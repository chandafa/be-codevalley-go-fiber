package main

import (
	"log"
	"time"

	"code-valley-api/internal/config"
	"code-valley-api/internal/database"
	"code-valley-api/internal/models"
	"code-valley-api/internal/utils"

	"github.com/google/uuid"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	if err := database.Initialize(cfg); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	db := database.GetDB()

	log.Println("Starting database seeding...")

	// Create admin user
	adminPassword, _ := utils.HashPassword("admin123")
	admin := &models.User{
		ID:           uuid.New(),
		Email:        "admin@codevalley.com",
		Username:     "admin",
		PasswordHash: adminPassword,
		Bio:          "Game Administrator",
		Role:         models.RoleAdmin,
		Level:        99,
		EXP:          9999,
		Coins:        10000,
	}
	db.FirstOrCreate(&admin, "email = ?", admin.Email)

	// Create sample player
	playerPassword, _ := utils.HashPassword("player123")
	player := &models.User{
		ID:           uuid.New(),
		Email:        "player@example.com",
		Username:     "coder123",
		PasswordHash: playerPassword,
		Bio:          "Aspiring developer",
		Role:         models.RolePlayer,
		Level:        1,
		EXP:          0,
		Coins:        100,
	}
	db.FirstOrCreate(&player, "email = ?", player.Email)

	// Create NPCs
	npcs := []models.NPC{
		{
			ID:        uuid.New(),
			Name:      "Marcus the Mentor",
			Role:      models.NPCRoleMentor,
			Dialogue:  "Welcome to Code Valley! I'm here to guide you on your programming journey.",
			Location:  "Town Center",
			AvatarURL: "https://example.com/marcus.jpg",
			IsActive:  true,
		},
		{
			ID:        uuid.New(),
			Name:      "Sarah the Client",
			Role:      models.NPCRoleClient,
			Dialogue:  "I have some projects that need a skilled developer. Are you up for the challenge?",
			Location:  "Business District",
			AvatarURL: "https://example.com/sarah.jpg",
			IsActive:  true,
		},
		{
			ID:        uuid.New(),
			Name:      "Bob the Villager",
			Role:      models.NPCRoleVillager,
			Dialogue:  "Life in Code Valley is great! Everyone here loves programming.",
			Location:  "Residential Area",
			AvatarURL: "https://example.com/bob.jpg",
			IsActive:  true,
		},
	}

	for _, npc := range npcs {
		db.FirstOrCreate(&npc, "name = ?", npc.Name)
	}

	// Create sample quests
	quests := []models.Quest{
		{
			ID:          uuid.New(),
			Title:       "First Steps in Programming",
			Description: "Learn the basics of programming by completing your first Hello World program.",
			RewardCoins: 50,
			RewardEXP:   100,
			RequiredItems: models.RequiredItems{
				"hello_world_code": 1,
			},
			IsRepeatable: false,
			IsActive:     true,
		},
		{
			ID:          uuid.New(),
			Title:       "Data Structure Challenge",
			Description: "Implement a basic data structure to prove your understanding.",
			RewardCoins: 100,
			RewardEXP:   200,
			RequiredItems: models.RequiredItems{
				"data_structure_implementation": 1,
				"test_cases": 3,
			},
			IsRepeatable: false,
			IsActive:     true,
		},
		{
			ID:          uuid.New(),
			Title:       "Daily Code Review",
			Description: "Review and improve existing code. This quest can be repeated daily.",
			RewardCoins: 25,
			RewardEXP:   50,
			RequiredItems: models.RequiredItems{
				"code_review": 1,
			},
			IsRepeatable: true,
			IsActive:     true,
		},
	}

	for _, quest := range quests {
		db.FirstOrCreate(&quest, "title = ?", quest.Title)
	}

	// Create achievements
	achievements := []models.Achievement{
		{
			ID:          uuid.New(),
			Title:       "First Steps",
			Description: "Complete your first quest",
			RewardCoins: 100,
			RewardEXP:   50,
			Conditions: models.Conditions{
				"quests_completed": 1,
			},
			IconURL:  "https://example.com/first_steps.png",
			IsActive: true,
		},
		{
			ID:          uuid.New(),
			Title:       "Level Up!",
			Description: "Reach level 5",
			RewardCoins: 250,
			RewardEXP:   100,
			Conditions: models.Conditions{
				"level": 5,
			},
			IconURL:  "https://example.com/level_up.png",
			IsActive: true,
		},
		{
			ID:          uuid.New(),
			Title:       "Code Master",
			Description: "Complete 10 quests",
			RewardCoins: 500,
			RewardEXP:   300,
			Conditions: models.Conditions{
				"quests_completed": 10,
			},
			IconURL:  "https://example.com/code_master.png",
			IsActive: true,
		},
	}

	for _, achievement := range achievements {
		db.FirstOrCreate(&achievement, "title = ?", achievement.Title)
	}

	// Create daily tasks
	today := time.Now()
	dailyTasks := []models.DailyTask{
		{
			ID:          uuid.New(),
			TaskName:    "Code 30 minutes",
			TaskType:    "coding",
			Description: "Spend at least 30 minutes coding today",
			RewardEXP:   25,
			RewardCoins: 10,
			Date:        today,
			IsActive:    true,
		},
		{
			ID:          uuid.New(),
			TaskName:    "Learn something new",
			TaskType:    "learning",
			Description: "Read about a new programming concept or technology",
			RewardEXP:   20,
			RewardCoins: 15,
			Date:        today,
			IsActive:    true,
		},
		{
			ID:          uuid.New(),
			TaskName:    "Help a fellow coder",
			TaskType:    "community",
			Description: "Answer a question or help someone in the community",
			RewardEXP:   30,
			RewardCoins: 20,
			Date:        today,
			IsActive:    true,
		},
	}

	for _, task := range dailyTasks {
		db.FirstOrCreate(&task, "task_name = ? AND date = ?", task.TaskName, task.Date)
	}

	// Create sample inventory items
	inventoryItems := []models.Inventory{
		{
			ID:       uuid.New(),
			UserID:   player.ID,
			ItemName: "Beginner's Keyboard",
			Quantity: 1,
			ItemType: models.ItemTypeTool,
		},
		{
			ID:       uuid.New(),
			UserID:   player.ID,
			ItemName: "JavaScript Snippet Collection",
			Quantity: 5,
			ItemType: models.ItemTypeSnippet,
		},
		{
			ID:       uuid.New(),
			UserID:   player.ID,
			ItemName: "Coffee Beans",
			Quantity: 10,
			ItemType: models.ItemTypeResource,
		},
	}

	for _, item := range inventoryItems {
		db.FirstOrCreate(&item, "user_id = ? AND item_name = ?", item.UserID, item.ItemName)
	}

	log.Println("Database seeding completed successfully!")
}