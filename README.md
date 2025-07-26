# Code Valley - RPG Life Simulation Game API

A complete REST API backend for "Code Valley" - an RPG life simulation game where players are junior programmers building careers in a coding village, inspired by Stardew Valley.

## 🚀 Features

- **User Management**: Registration, authentication, profile management
- **Quest System**: Create, start, complete quests with rewards
- **Inventory System**: Manage tools, code snippets, and resources
- **NPC Interactions**: Meet mentors, clients, and villagers
- **Daily Tasks**: Complete daily coding challenges
- **Achievement System**: Unlock achievements based on progress
- **Story Progress**: Track player's journey through chapters
- **Code Battles**: Participate in coding challenges
- **Admin Panel**: Manage game content and users

## 🛠 Technical Stack

- **Framework**: Go with Fiber web framework
- **Database**: MySQL with GORM ORM
- **Authentication**: JWT with bcrypt password hashing
- **Architecture**: Clean architecture with layered design
- **Middleware**: CORS, Rate limiting, Logging, Error handling
- **Validation**: Struct validation with custom validators
- **Configuration**: Environment-based with .env support

## 📁 Project Structure

```
code-valley-api/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── database/        # Database connection and setup
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models and structs
│   ├── repositories/    # Data access layer
│   ├── services/        # Business logic layer
│   ├── routes/          # Route definitions
│   └── utils/           # Utility functions
├── seeders/             # Database seeders
├── .env.example         # Environment variables template
├── go.mod              # Go module definition
└── README.md           # This file
```

## 🚀 Getting Started

### Prerequisites

- Go 1.21 or higher
- MySQL 8.0 or higher
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd code-valley-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials and configuration
   ```

4. **Create MySQL database**
   ```sql
   CREATE DATABASE code_valley;
   ```

5. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

6. **Seed the database (optional)**
   ```bash
   go run seeders/seed.go
   ```

### Environment Variables

```env
DB_HOST=pongo.kencang.com
DB_PORT=3306
DB_USER=academyc_root_pp
DB_PASSWORD=Langkahpemula123
DB_NAME=academyc_code_valley
PORT=8000

JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRE_HOURS=24

CORS_ORIGIN=*

RATE_LIMIT_MAX=100
RATE_LIMIT_EXPIRATION=1

LOG_LEVEL=info
```

## 📚 API Documentation

### Base URL
```
http://localhost:8000/api/v1
```

### Response Format
All API responses follow this structure:
```json
{
  "success": boolean,
  "message": "descriptive message",
  "data": any | null
}
```

### Authentication Endpoints

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "username",
  "password": "password123"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

#### Get Profile (Protected)
```http
GET /api/v1/auth/me
Authorization: Bearer <jwt-token>
```

#### Update Profile (Protected)
```http
PUT /api/v1/auth/profile
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "username": "new_username",
  "bio": "Updated bio",
  "avatar_url": "https://example.com/avatar.jpg"
}
```

### Quest Endpoints

#### Get All Quests (Protected)
```http
GET /api/v1/quests?page=1&per_page=10
Authorization: Bearer <jwt-token>
```

#### Get Quest Details (Protected)
```http
GET /api/v1/quests/:id
Authorization: Bearer <jwt-token>
```

#### Start Quest (Protected)
```http
POST /api/v1/quests/:id/start
Authorization: Bearer <jwt-token>
```

#### Complete Quest (Protected)
```http
POST /api/v1/quests/:id/complete
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "submitted_items": {
    "item_name": quantity
  }
}
```

#### Get User Progress (Protected)
```http
GET /api/v1/quests/progress
Authorization: Bearer <jwt-token>
```

### Admin Endpoints (Admin Role Required)

#### Create Quest
```http
POST /api/v1/admin/quests
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json

{
  "title": "Quest Title",
  "description": "Quest description",
  "reward_coins": 100,
  "reward_exp": 50,
  "required_items": {
    "item_name": 2
  },
  "is_repeatable": false,
  "is_active": true
}
```

#### Update Quest
```http
PUT /api/v1/admin/quests/:id
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json

{
  "title": "Updated Quest Title",
  "description": "Updated description",
  // ... other fields
}
```

#### Delete Quest
```http
DELETE /api/v1/admin/quests/:id
Authorization: Bearer <admin-jwt-token>
```

## 🗄 Database Schema

### Core Models

- **User**: Player accounts with authentication and game stats
- **Inventory**: Player items (tools, code snippets, resources)
- **Quest**: Available quests with requirements and rewards
- **UserQuestProgress**: Player progress on specific quests
- **NPC**: Non-player characters with dialogue and locations
- **DailyTask**: Daily challenges for players
- **UserDailyTaskProgress**: Player completion of daily tasks
- **Achievement**: Unlockable achievements with conditions
- **UserAchievement**: Player's unlocked achievements
- **StoryProgress**: Player's story chapter progress
- **CodeBattle**: Coding challenges and battles

## 🔧 Middleware

- **Authentication**: JWT token validation
- **CORS**: Cross-origin request handling
- **Rate Limiting**: Request rate limiting per IP
- **Logging**: Request/response logging
- **Error Handling**: Centralized error handling
- **Recovery**: Panic recovery

## 🧪 Full Routes.go
func SetupRoutes(app *fiber.App, cfg *config.Config) {
	// Initialize services
	authService := services.NewAuthService(cfg)
	questService := services.NewQuestService()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	questHandler := handlers.NewQuestHandler(questService)

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

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"message": "Code Valley API is running",
		})
	})
}


```

## 📈 Performance Tips

1. **Database Indexing**: Add indexes on frequently queried columns
2. **Connection Pooling**: Configure appropriate connection pool settings
3. **Caching**: Implement Redis for session management and caching
4. **Pagination**: Always use pagination for list endpoints
5. **Rate Limiting**: Implement appropriate rate limits


prompt frontend
1. pisahkan dashboard admin dan player
2. buatkan halaman npc, progres, invectory, portfolio, map dan profile, serta halaman lain yang dibutuhkan
3. sesuaikan websocket dengan backend supaya realtime tanpa refresh
4. perbaiki responsive halaman pada layar mobile
5, sesuaikan fitur dan routes api nya dengan dokumenasi backend berikut :

promp backend 
1. buatkan websocket realtime tanpa refresh
2. perbaiki query tunning dan normalisasi
3. Route Tambahan :
1. Friend System / Social
Biar player bisa saling bantu atau kompetitif:

http
```
Edit
GET    /api/v1/friends                       # List teman
POST   /api/v1/friends/:username/add         # Kirim permintaan pertemanan
POST   /api/v1/friends/:username/accept      # Terima permintaan
DELETE /api/v1/friends/:username/remove      # Hapus teman
GET    /api/v1/friends/online                # Teman yang sedang online
Gunanya untuk leaderboard, bantuan quest, atau PvP code battle.

2. Leaderboard / Ranking
Fitur kompetitif ringan:

http
```
Edit
GET /api/v1/leaderboard/coins       # Top player by coins
GET /api/v1/leaderboard/exp         # Top player by EXP
GET /api/v1/leaderboard/tasks       # Top by daily task completion
3. Shop & Economy
Biar hasil kerja coding bisa dibelanjakan.

http
```
Edit
GET    /api/v1/shop/items                 # Lihat item toko
POST   /api/v1/shop/items/:id/buy         # Beli item
POST   /api/v1/shop/items/:id/sell        # Jual item dari inventory
Bisa integrasi dengan inventory dan coin system.

4. Mailbox / Notifications
Buat feedback event atau info penting dari NPC/mentor:

http
```
Edit
GET    /api/v1/notifications              # Semua notifikasi
POST   /api/v1/notifications/mark-read   # Tandai sebagai dibaca
5. Mentorship / NPC Relationship
Interaksi sosial ala Stardew Valley.

http
```
Edit
GET    /api/v1/npc/:id/profile            # Lihat profil NPC
POST   /api/v1/npc/:id/interact           # Interaksi (berbincang, beri hadiah)
GET    /api/v1/npc/:id/relationship       # Level kedekatan dengan NPC
Semakin akrab, bisa dapat quest eksklusif atau unlock skill.

6. Skill Tree / Upgrade
Biar player merasa progres:

http
```
Edit
GET    /api/v1/skills                      # Lihat semua skill
POST   /api/v1/skills/:id/unlock           # Unlock/upgrade skill
Misalnya: "Debugging +10%", "Faster Compile Time", dll.

7. Mini Games / Coding Challenges
Untuk variasi gameplay:

http
```
Edit
GET    /api/v1/minigames/list             # Lihat list mini games
POST   /api/v1/minigames/:id/start        # Mulai game
POST   /api/v1/minigames/:id/submit       # Kirim hasil
Bisa coding quiz, logic puzzle, regex challenge.

8. Item Crafting / Code Combining
Seperti gabung script jadi alat baru.

http
```
Edit
GET    /api/v1/crafting/recipes           # Semua resep
POST   /api/v1/crafting/execute           # Craft item
9. Daily Login Reward / Calendar System
Supaya player rajin main tiap hari.

http
```
Edit
GET    /api/v1/rewards/daily              # Cek status harian
POST   /api/v1/rewards/daily/claim        # Ambil reward harian
10. Event System / Time-based Missions
Mirip festival di Stardew atau "Hackathon Week".

http
```
Edit
GET    /api/v1/events/active              # Event aktif sekarang
GET    /api/v1/events/:id                 # Detail event
POST   /api/v1/events/:id/join            # Gabung event


Buatkan Fitur juga fitur berikut :
api/v1/badges	Visual penghargaan, mirip achievement
api/v1/guilds	Buat guild/komunitas programmer, unlock fitur kooperatif
api/v1/marketplace	Jual beli script antar user
api/v1/tutorials	Panduan coding, bisa unlock setelah quest tertentu
api/v1/statistics	Statistik performa user (waktu main, progress, dsb)