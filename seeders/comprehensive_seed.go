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

	log.Println("Starting comprehensive database seeding...")

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

	// Create sample players
	players := []models.User{
		{
			ID:           uuid.New(),
			Email:        "alice@example.com",
			Username:     "alice_codes",
			PasswordHash: adminPassword,
			Bio:          "Frontend enthusiast",
			Role:         models.RolePlayer,
			Level:        5,
			EXP:          450,
			Coins:        500,
		},
		{
			ID:           uuid.New(),
			Email:        "bob@example.com",
			Username:     "bob_backend",
			PasswordHash: adminPassword,
			Bio:          "Backend wizard",
			Role:         models.RolePlayer,
			Level:        7,
			EXP:          680,
			Coins:        750,
		},
		{
			ID:           uuid.New(),
			Email:        "charlie@example.com",
			Username:     "charlie_fullstack",
			PasswordHash: adminPassword,
			Bio:          "Full-stack developer",
			Role:         models.RolePlayer,
			Level:        3,
			EXP:          280,
			Coins:        300,
		},
	}

	for _, player := range players {
		db.FirstOrCreate(&player, "email = ?", player.Email)
	}

	// Create Maps
	maps := []models.Map{
		{
			ID:     uuid.New(),
			Name:   "village",
			Type:   models.MapTypeVillage,
			Width:  50,
			Height: 50,
			Layout: models.MapLayout{
				"spawn_point": map[string]int{"x": 25, "y": 25},
				"buildings":   []map[string]interface{}{
					{"type": "town_hall", "x": 25, "y": 20},
					{"type": "shop", "x": 30, "y": 25},
					{"type": "library", "x": 20, "y": 25},
				},
			},
			Description: "The main village where programmers gather",
			IsActive:    true,
		},
		{
			ID:     uuid.New(),
			Name:   "code_mine",
			Type:   models.MapTypeCodeMine,
			Width:  40,
			Height: 40,
			Layout: models.MapLayout{
				"entrance": map[string]int{"x": 20, "y": 35},
				"levels":   5,
			},
			Description: "Deep caves where you can mine for algorithms and data structures",
			IsActive:    true,
		},
		{
			ID:     uuid.New(),
			Name:   "data_farm",
			Type:   models.MapTypeDataFarm,
			Width:  60,
			Height: 40,
			Layout: models.MapLayout{
				"plots": 100,
				"greenhouse": map[string]int{"x": 30, "y": 20},
			},
			Description: "Your personal coding farm where you grow and nurture code",
			IsActive:    true,
		},
	}

	for _, mapData := range maps {
		db.FirstOrCreate(&mapData, "name = ?", mapData.Name)
	}

	// Create NPCs
	npcs := []models.NPC{
		{
			ID:        uuid.New(),
			Name:      "Marcus the Mentor",
			Role:      models.NPCRoleMentor,
			Dialogue:  "Welcome to Code Valley! I'm here to guide you on your programming journey.",
			Location:  "Village Center",
			AvatarURL: "https://images.pexels.com/photos/220453/pexels-photo-220453.jpeg?auto=compress&cs=tinysrgb&w=200",
			IsActive:  true,
		},
		{
			ID:        uuid.New(),
			Name:      "Sarah the Client",
			Role:      models.NPCRoleClient,
			Dialogue:  "I have some projects that need a skilled developer. Are you up for the challenge?",
			Location:  "Business District",
			AvatarURL: "https://images.pexels.com/photos/415829/pexels-photo-415829.jpeg?auto=compress&cs=tinysrgb&w=200",
			IsActive:  true,
		},
		{
			ID:        uuid.New(),
			Name:      "Bob the Villager",
			Role:      models.NPCRoleVillager,
			Dialogue:  "Life in Code Valley is great! Everyone here loves programming.",
			Location:  "Residential Area",
			AvatarURL: "https://images.pexels.com/photos/614810/pexels-photo-614810.jpeg?auto=compress&cs=tinysrgb&w=200",
			IsActive:  true,
		},
		{
			ID:        uuid.New(),
			Name:      "Dr. Debug",
			Role:      models.NPCRoleMentor,
			Dialogue:  "Bugs are not your enemy - they're learning opportunities!",
			Location:  "Code Mine Entrance",
			AvatarURL: "https://images.pexels.com/photos/2182970/pexels-photo-2182970.jpeg?auto=compress&cs=tinysrgb&w=200",
			IsActive:  true,
		},
	}

	for _, npc := range npcs {
		db.FirstOrCreate(&npc, "name = ?", npc.Name)
	}

	// Create Quests
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
			Title:       "Bug Hunt",
			Description: "Find and fix 5 bugs in the legacy codebase.",
			RewardCoins: 75,
			RewardEXP:   150,
			RequiredItems: models.RequiredItems{
				"bug_fixes": 5,
			},
			IsRepeatable: true,
			IsActive:     true,
		},
		{
			ID:          uuid.New(),
			Title:       "Code Review Master",
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

	// Create Shop Items
	shopItems := []models.ShopItem{
		{
			ID:          uuid.New(),
			Name:        "Basic Keyboard",
			Description: "A reliable keyboard for coding",
			Price:       100,
			ItemType:    models.ShopItemTypeTool,
			IconURL:     "https://images.pexels.com/photos/2115257/pexels-photo-2115257.jpeg?auto=compress&cs=tinysrgb&w=200",
			IsAvailable: true,
			Stock:       10,
		},
		{
			ID:          uuid.New(),
			Name:        "Debug Magnifier",
			Description: "Helps you spot bugs more easily",
			Price:       250,
			ItemType:    models.ShopItemTypeTool,
			IconURL:     "https://images.pexels.com/photos/1181675/pexels-photo-1181675.jpeg?auto=compress&cs=tinysrgb&w=200",
			IsAvailable: true,
			Stock:       5,
		},
		{
			ID:          uuid.New(),
			Name:        "Coffee Beans",
			Description: "Essential fuel for programmers",
			Price:       10,
			ItemType:    models.ShopItemTypeResource,
			IconURL:     "https://images.pexels.com/photos/894695/pexels-photo-894695.jpeg?auto=compress&cs=tinysrgb&w=200",
			IsAvailable: true,
			Stock:       -1, // Unlimited
		},
		{
			ID:          uuid.New(),
			Name:        "Algorithm Seeds",
			Description: "Plant these to grow powerful algorithms",
			Price:       50,
			ItemType:    models.ShopItemTypeResource,
			IconURL:     "https://images.pexels.com/photos/1459505/pexels-photo-1459505.jpeg?auto=compress&cs=tinysrgb&w=200",
			IsAvailable: true,
			Stock:       20,
		},
	}

	for _, item := range shopItems {
		db.FirstOrCreate(&item, "name = ?", item.Name)
	}

	// Create Achievements
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
			IconURL:  "https://images.pexels.com/photos/1181675/pexels-photo-1181675.jpeg?auto=compress&cs=tinysrgb&w=200",
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
			IconURL:  "https://images.pexels.com/photos/1181675/pexels-photo-1181675.jpeg?auto=compress&cs=tinysrgb&w=200",
			IsActive: true,
		},
		{
			ID:          uuid.New(),
			Title:       "Bug Squasher",
			Description: "Fix 50 bugs",
			RewardCoins: 500,
			RewardEXP:   300,
			Conditions: models.Conditions{
				"bugs_fixed": 50,
			},
			IconURL:  "https://images.pexels.com/photos/1181675/pexels-photo-1181675.jpeg?auto=compress&cs=tinysrgb&w=200",
			IsActive: true,
		},
	}

	for _, achievement := range achievements {
		db.FirstOrCreate(&achievement, "title = ?", achievement.Title)
	}

	// Create Skills
	skills := []models.Skill{
		{
			ID:          uuid.New(),
			Name:        "Debugging Mastery",
			Description: "Improve your debugging efficiency",
			Category:    models.SkillCategoryDebugging,
			MaxLevel:    10,
			BaseCost:    100,
			Effects: models.SkillEffects{
				"debug_speed": "+10% per level",
				"bug_detection": "+5% per level",
			},
			IconURL:  "https://images.pexels.com/photos/1181675/pexels-photo-1181675.jpeg?auto=compress&cs=tinysrgb&w=200",
			IsActive: true,
		},
		{
			ID:          uuid.New(),
			Name:        "Code Optimization",
			Description: "Write more efficient code",
			Category:    models.SkillCategoryOptimization,
			MaxLevel:    10,
			BaseCost:    150,
			Effects: models.SkillEffects{
				"performance": "+15% per level",
				"memory_usage": "-5% per level",
			},
			IconURL:  "https://images.pexels.com/photos/1181675/pexels-photo-1181675.jpeg?auto=compress&cs=tinysrgb&w=200",
			IsActive: true,
		},
		{
			ID:          uuid.New(),
			Name:        "Social Coding",
			Description: "Better collaboration with other developers",
			Category:    models.SkillCategorySocial,
			MaxLevel:    5,
			BaseCost:    200,
			Effects: models.SkillEffects{
				"team_bonus": "+20% per level",
				"code_review_speed": "+10% per level",
			},
			IconURL:  "https://images.pexels.com/photos/1181675/pexels-photo-1181675.jpeg?auto=compress&cs=tinysrgb&w=200",
			IsActive: true,
		},
	}

	for _, skill := range skills {
		db.FirstOrCreate(&skill, "name = ?", skill.Name)
	}

	// Create Mini Games
	miniGames := []models.MiniGame{
		{
			ID:          uuid.New(),
			Name:        "JavaScript Quiz",
			Description: "Test your JavaScript knowledge",
			Type:        models.MiniGameTypeQuiz,
			Difficulty:  models.DifficultyEasy,
			Config: models.GameConfig{
				"questions": 10,
				"time_limit": 300,
				"topics": []string{"variables", "functions", "arrays"},
			},
			RewardCoins: 50,
			RewardEXP:   25,
			TimeLimit:   300,
			IsActive:    true,
		},
		{
			ID:          uuid.New(),
			Name:        "Algorithm Puzzle",
			Description: "Solve algorithmic challenges",
			Type:        models.MiniGameTypeAlgorithm,
			Difficulty:  models.DifficultyMedium,
			Config: models.GameConfig{
				"problems": 5,
				"difficulty": "medium",
				"topics": []string{"sorting", "searching", "dynamic_programming"},
			},
			RewardCoins: 100,
			RewardEXP:   75,
			TimeLimit:   600,
			IsActive:    true,
		},
		{
			ID:          uuid.New(),
			Name:        "Regex Master",
			Description: "Master regular expressions",
			Type:        models.MiniGameTypeRegex,
			Difficulty:  models.DifficultyHard,
			Config: models.GameConfig{
				"patterns": 8,
				"complexity": "high",
			},
			RewardCoins: 150,
			RewardEXP:   100,
			TimeLimit:   450,
			IsActive:    true,
		},
	}

	for _, game := range miniGames {
		db.FirstOrCreate(&game, "name = ?", game.Name)
	}

	// Create Guilds
	guilds := []models.Guild{
		{
			ID:          uuid.New(),
			Name:        "Code Warriors",
			Description: "Elite programmers unite for epic coding battles!",
			OwnerID:     players[0].ID,
			MaxMembers:  20,
			Level:       3,
			EXP:         450,
			IconURL:     "https://images.pexels.com/photos/1181675/pexels-photo-1181675.jpeg?auto=compress&cs=tinysrgb&w=200",
			IsPublic:    true,
		},
		{
			ID:          uuid.New(),
			Name:        "Debug Squad",
			Description: "We hunt bugs together and make code better",
			OwnerID:     players[1].ID,
			MaxMembers:  15,
			Level:       2,
			EXP:         200,
			IconURL:     "https://images.pexels.com/photos/1181675/pexels-photo-1181675.jpeg?auto=compress&cs=tinysrgb&w=200",
			IsPublic:    true,
		},
	}

	for _, guild := range guilds {
		db.FirstOrCreate(&guild, "name = ?", guild.Name)
	}

	// Create Events
	events := []models.Event{
		{
			ID:          uuid.New(),
			Name:        "Spring Hackathon",
			Description: "A week-long coding competition with amazing prizes!",
			Type:        models.EventTypeHackathon,
			StartDate:   time.Now().Add(24 * time.Hour),
			EndDate:     time.Now().Add(8 * 24 * time.Hour),
			Requirements: models.EventRewards{
				"min_level": 3,
				"entry_fee": 100,
			},
			Rewards: models.EventRewards{
				"first_place":  map[string]interface{}{"coins": 1000, "exp": 500, "badge": "Hackathon Champion"},
				"second_place": map[string]interface{}{"coins": 500, "exp": 250, "badge": "Hackathon Runner-up"},
				"third_place":  map[string]interface{}{"coins": 250, "exp": 125, "badge": "Hackathon Participant"},
			},
			MaxParticipants: 50,
			IsActive:        true,
		},
		{
			ID:          uuid.New(),
			Name:        "Bug Squashing Festival",
			Description: "Help clean up the codebase and earn rewards!",
			Type:        models.EventTypeFestival,
			StartDate:   time.Now().Add(3 * 24 * time.Hour),
			EndDate:     time.Now().Add(10 * 24 * time.Hour),
			Requirements: models.EventRewards{
				"min_level": 1,
			},
			Rewards: models.EventRewards{
				"bugs_fixed_1":  map[string]interface{}{"coins": 50, "exp": 25},
				"bugs_fixed_5":  map[string]interface{}{"coins": 200, "exp": 100},
				"bugs_fixed_10": map[string]interface{}{"coins": 500, "exp": 250, "badge": "Bug Hunter"},
			},
			MaxParticipants: -1,
			IsActive:        true,
		},
	}

	for _, event := range events {
		db.FirstOrCreate(&event, "name = ?", event.Name)
	}

	// Create Daily Tasks
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
		{
			ID:          uuid.New(),
			TaskName:    "Fix a bug",
			TaskType:    "debugging",
			Description: "Find and fix at least one bug in any codebase",
			RewardEXP:   35,
			RewardCoins: 25,
			Date:        today,
			IsActive:    true,
		},
	}

	for _, task := range dailyTasks {
		db.FirstOrCreate(&task, "task_name = ? AND date = ?", task.TaskName, task.Date)
	}

	// Create sample inventory items for players
	inventoryItems := []models.Inventory{
		{
			ID:       uuid.New(),
			UserID:   players[0].ID,
			ItemName: "Beginner's Keyboard",
			Quantity: 1,
			ItemType: models.ItemTypeTool,
		},
		{
			ID:       uuid.New(),
			UserID:   players[0].ID,
			ItemName: "JavaScript Snippet Collection",
			Quantity: 5,
			ItemType: models.ItemTypeSnippet,
		},
		{
			ID:       uuid.New(),
			UserID:   players[1].ID,
			ItemName: "Coffee Beans",
			Quantity: 10,
			ItemType: models.ItemTypeResource,
		},
		{
			ID:       uuid.New(),
			UserID:   players[1].ID,
			ItemName: "Debug Tool",
			Quantity: 2,
			ItemType: models.ItemTypeTool,
		},
	}

	for _, item := range inventoryItems {
		db.FirstOrCreate(&item, "user_id = ? AND item_name = ?", item.UserID, item.ItemName)
	}

	// Create friendships
	friendships := []models.Friendship{
		{
			ID:          uuid.New(),
			RequesterID: players[0].ID,
			AddresseeID: players[1].ID,
			Status:      models.FriendshipStatusAccepted,
		},
		{
			ID:          uuid.New(),
			RequesterID: players[1].ID,
			AddresseeID: players[2].ID,
			Status:      models.FriendshipStatusPending,
		},
	}

	for _, friendship := range friendships {
		db.FirstOrCreate(&friendship, "requester_id = ? AND addressee_id = ?", friendship.RequesterID, friendship.AddresseeID)
	}

	// Create World Objects
	villageMap := maps[0] // Village map
	worldObjects := []models.WorldObject{
		{
			ID:         uuid.New(),
			MapID:      villageMap.ID,
			ObjectType: models.ObjectTypeTree,
			PosX:       10,
			PosY:       15,
			State: models.ObjectState{
				"hp":      3,
				"type":    "oak",
				"matured": true,
			},
			IsActive: true,
		},
		{
			ID:         uuid.New(),
			MapID:      villageMap.ID,
			ObjectType: models.ObjectTypeRock,
			PosX:       35,
			PosY:       20,
			State: models.ObjectState{
				"hp":   2,
				"type": "granite",
			},
			IsActive: true,
		},
		{
			ID:         uuid.New(),
			MapID:      villageMap.ID,
			ObjectType: models.ObjectTypeChest,
			PosX:       40,
			PosY:       30,
			State: models.ObjectState{
				"is_looted": false,
				"contents":  []string{"Debug Tool", "Coffee Beans"},
			},
			IsActive: true,
		},
	}

	for _, obj := range worldObjects {
		db.FirstOrCreate(&obj, "map_id = ? AND pos_x = ? AND pos_y = ?", obj.MapID, obj.PosX, obj.PosY)
	}

	// Create Player Positions
	playerPositions := []models.PlayerPosition{
		{
			ID:        uuid.New(),
			UserID:    players[0].ID,
			MapID:     villageMap.ID,
			PosX:      25,
			PosY:      25,
			Direction: "down",
			LastMoved: time.Now(),
		},
		{
			ID:        uuid.New(),
			UserID:    players[1].ID,
			MapID:     villageMap.ID,
			PosX:      26,
			PosY:      25,
			Direction: "left",
			LastMoved: time.Now(),
		},
	}

	for _, pos := range playerPositions {
		db.FirstOrCreate(&pos, "user_id = ?", pos.UserID)
	}

	// Create NPC Positions
	npcPositions := []models.NPCPosition{
		{
			ID:        uuid.New(),
			NPCID:     npcs[0].ID, // Marcus the Mentor
			MapID:     villageMap.ID,
			PosX:      25,
			PosY:      20,
			Direction: "down",
		},
		{
			ID:        uuid.New(),
			NPCID:     npcs[1].ID, // Sarah the Client
			MapID:     villageMap.ID,
			PosX:      30,
			PosY:      25,
			Direction: "left",
		},
	}

	for _, pos := range npcPositions {
		db.FirstOrCreate(&pos, "npc_id = ?", pos.NPCID)
	}

	// Create Game Clock
	gameClock := &models.GameClock{
		ID:         uuid.New(),
		GameYear:   1,
		GameSeason: "spring",
		GameDay:    1,
		GameHour:   6,
		GameMinute: 0,
		IsPaused:   false,
		TimeScale:  1.0,
	}
	db.FirstOrCreate(&gameClock)

	// Create Daily Rewards
	dailyRewards := []models.DailyReward{
		{
			ID:          uuid.New(),
			Day:         1,
			RewardCoins: 50,
			RewardEXP:   25,
			BonusItem:   "Coffee Beans",
			IsActive:    true,
		},
		{
			ID:          uuid.New(),
			Day:         2,
			RewardCoins: 75,
			RewardEXP:   35,
			BonusItem:   "Debug Tool",
			IsActive:    true,
		},
		{
			ID:          uuid.New(),
			Day:         7,
			RewardCoins: 200,
			RewardEXP:   100,
			BonusItem:   "Premium Keyboard",
			IsActive:    true,
		},
	}

	for _, reward := range dailyRewards {
		db.FirstOrCreate(&reward, "day = ?", reward.Day)
	}

	log.Println("Comprehensive database seeding completed successfully!")
	log.Println("Created:")
	log.Println("- 1 Admin user and 3 sample players")
	log.Println("- 3 Maps (Village, Code Mine, Data Farm)")
	log.Println("- 4 NPCs with different roles")
	log.Println("- 4 Quests with various difficulties")
	log.Println("- 4 Shop items")
	log.Println("- 3 Achievements")
	log.Println("- 3 Skills")
	log.Println("- 3 Mini games")
	log.Println("- 2 Guilds")
	log.Println("- 2 Events")
	log.Println("- 4 Daily tasks")
	log.Println("- Sample inventory items")
	log.Println("- Friend relationships")
	log.Println("- World objects and positions")
	log.Println("- Game clock and daily rewards")
}