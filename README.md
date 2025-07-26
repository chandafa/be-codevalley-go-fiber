# Code Valley - RPG Life Simulation Game API

A complete REST API backend for "Code Valley" - an RPG life simulation game where players are junior programmers building careers in a coding village, inspired by Stardew Valley.

## 🚀 Features

- **User Management**: Registration, authentication, profile management, avatar upload
- **Real-time WebSocket**: Live updates for quests, friends, achievements, and more
- **Quest System**: Create, start, complete quests with rewards
- **Friend System**: Add friends, see online status, social interactions
- **Inventory System**: Manage tools, code snippets, and resources
- **NPC Interactions**: Meet mentors, clients, and villagers with relationship levels
- **Daily Tasks**: Complete daily coding challenges
- **Achievement & Badge System**: Unlock achievements and collect badges
- **Skill Tree**: Upgrade programming skills and abilities
- **Story Progress**: Track player's journey through chapters
- **Mini Games**: Coding challenges, puzzles, and brain teasers
- **Crafting System**: Combine items to create new tools
- **Shop & Economy**: Buy and sell items with coins
- **Guild System**: Create and join programming communities
- **Event System**: Participate in time-based missions and festivals
- **Marketplace**: Trade items with other players
- **Tutorial System**: Learn new programming concepts
- **Statistics**: Comprehensive player analytics
- **Admin Panel**: Manage game content, users, and analytics

## 🛠 Technical Stack

- **Framework**: Go with Fiber web framework
- **Database**: MySQL with GORM ORM
- **WebSocket**: Real-time communication with gorilla/websocket
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
│   ├── websocket/       # WebSocket hub and client management
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
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=code_valley
PORT=8000

JWT_SECRET=your-super-secret-jwt-key
JWT_EXPIRE_HOURS=24

CORS_ORIGIN=*

RATE_LIMIT_MAX=100
RATE_LIMIT_EXPIRATION=1

LOG_LEVEL=info
```

## 🌐 WebSocket Connection

Connect to WebSocket for real-time updates:
```
ws://localhost:8000/ws?token=YOUR_JWT_TOKEN
```

### WebSocket Events

#### Outgoing Events (Server → Client)
- `quest_update`: Quest progress changes
- `friend_request`: Friend system notifications
- `achievement_unlocked`: New achievements earned
- `level_up`: Level progression updates
- `user_status`: Friend online/offline status
- `event_broadcast`: Global announcements
- `dm_message`: Direct messages
- `guild_invitation`: Guild invitations
- `notification`: General notifications

#### Incoming Events (Client → Server)
- `ping`: Keep connection alive
- `dm_message`: Send direct message
- `dm_typing`: Typing indicator
- `quest_update`: Quest progress update

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

---

## 🔐 Authentication Endpoints

### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "username": "username",
  "password": "password123"
}
```

### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### Logout
```http
POST /api/v1/auth/logout
Authorization: Bearer <jwt-token>
```

### Get Profile
```http
GET /api/v1/auth/me
Authorization: Bearer <jwt-token>
```

### Update Profile
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

### Upload Avatar
```http
POST /api/v1/auth/avatar
Authorization: Bearer <jwt-token>
Content-Type: multipart/form-data

avatar: <file>
```

### Delete Account
```http
DELETE /api/v1/auth/delete
Authorization: Bearer <jwt-token>
```

---

## 🎯 Quest System

### Get All Quests
```http
GET /api/v1/quests?page=1&per_page=10
Authorization: Bearer <jwt-token>
```

### Get Quest Details
```http
GET /api/v1/quests/:id
Authorization: Bearer <jwt-token>
```

### Start Quest
```http
POST /api/v1/quests/:id/start
Authorization: Bearer <jwt-token>
```

### Complete Quest
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

### Get User Progress
```http
GET /api/v1/quests/progress
Authorization: Bearer <jwt-token>
```

---

## 👥 Friend System

### Get Friends List
```http
GET /api/v1/friends
Authorization: Bearer <jwt-token>
```

### Send Friend Request
```http
POST /api/v1/friends/:username/add
Authorization: Bearer <jwt-token>
```

### Accept Friend Request
```http
POST /api/v1/friends/:username/accept
Authorization: Bearer <jwt-token>
```

### Remove Friend
```http
DELETE /api/v1/friends/:username/remove
Authorization: Bearer <jwt-token>
```

### Get Online Friends
```http
GET /api/v1/friends/online
Authorization: Bearer <jwt-token>
```

---

## 🏆 Leaderboard System

### Top Players by Coins
```http
GET /api/v1/leaderboard/coins
```

### Top Players by EXP
```http
GET /api/v1/leaderboard/exp
```

### Top Players by Tasks Completed
```http
GET /api/v1/leaderboard/tasks
```

---

## 🛒 Shop & Economy

### Get Shop Items
```http
GET /api/v1/shop/items?page=1&per_page=10
Authorization: Bearer <jwt-token>
```

### Buy Item
```http
POST /api/v1/shop/items/:id/buy
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "quantity": 1
}
```

### Sell Item
```http
POST /api/v1/shop/items/:id/sell
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "quantity": 1
}
```

---

## 📦 Inventory Management

### Get User Inventory
```http
GET /api/v1/inventory
Authorization: Bearer <jwt-token>
```

### Get Specific Item
```http
GET /api/v1/inventory/:id
Authorization: Bearer <jwt-token>
```

### Use Item
```http
POST /api/v1/inventory/use/:id
Authorization: Bearer <jwt-token>
```

### Equip Item
```http
POST /api/v1/inventory/equip/:id
Authorization: Bearer <jwt-token>
```

### Unequip Item
```http
POST /api/v1/inventory/unequip/:id
Authorization: Bearer <jwt-token>
```

---

## 🔔 Notifications

### Get All Notifications
```http
GET /api/v1/notifications
Authorization: Bearer <jwt-token>
```

### Mark Notification as Read
```http
POST /api/v1/notifications/:id/read
Authorization: Bearer <jwt-token>
```

### Mark All as Read
```http
POST /api/v1/notifications/mark-read
Authorization: Bearer <jwt-token>
```

---

## 💬 Messaging System

### Send Direct Message
```http
POST /api/v1/messages/:username/send
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "message": "Hello there!"
}
```

### Get Conversations
```http
GET /api/v1/messages/conversations
Authorization: Bearer <jwt-token>
```

---

## 🏆 Achievements & Badges

### Get User Achievements
```http
GET /api/v1/achievements
Authorization: Bearer <jwt-token>
```

### Get User Badges
```http
GET /api/v1/badges
Authorization: Bearer <jwt-token>
```

---

## 🧠 Skill System

### Get All Skills
```http
GET /api/v1/skills
Authorization: Bearer <jwt-token>
```

### Upgrade Skill
```http
POST /api/v1/skills/:id/upgrade
Authorization: Bearer <jwt-token>
```

---

## 🏰 Guild System

### Create Guild
```http
POST /api/v1/guilds
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "name": "Code Warriors",
  "description": "Elite programmers unite!",
  "is_public": true
}
```

### Get All Guilds
```http
GET /api/v1/guilds?page=1&per_page=10
```

### Get Guild Details
```http
GET /api/v1/guilds/:id
Authorization: Bearer <jwt-token>
```

### Join Guild
```http
POST /api/v1/guilds/:id/join
Authorization: Bearer <jwt-token>
```

### Leave Guild
```http
POST /api/v1/guilds/:id/leave
Authorization: Bearer <jwt-token>
```

### Invite to Guild
```http
POST /api/v1/guilds/:id/invite/:username
Authorization: Bearer <jwt-token>
```

### Kick from Guild
```http
POST /api/v1/guilds/:id/kick/:username
Authorization: Bearer <jwt-token>
```

---

## 🎪 Event System

### Get Active Events
```http
GET /api/v1/events
```

### Get Event Details
```http
GET /api/v1/events/:id
Authorization: Bearer <jwt-token>
```

### Join Event
```http
POST /api/v1/events/:id/join
Authorization: Bearer <jwt-token>
```

### Complete Event
```http
POST /api/v1/events/:id/complete
Authorization: Bearer <jwt-token>
```

---

## 🧪 Crafting System

### Get Crafting Recipes
```http
GET /api/v1/crafting/recipes
Authorization: Bearer <jwt-token>
```

### Execute Crafting
```http
POST /api/v1/crafting/:id/execute
Authorization: Bearer <jwt-token>
```

---

## 🧩 Mini Games

### Get Mini Games List
```http
GET /api/v1/minigames
Authorization: Bearer <jwt-token>
```

### Start Mini Game
```http
POST /api/v1/minigames/:id/start
Authorization: Bearer <jwt-token>
```

### Submit Mini Game Result
```http
POST /api/v1/minigames/:id/submit
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "answers": ["answer1", "answer2"],
  "score": 85,
  "time_taken": 120
}
```

---

## 🎁 Daily Rewards

### Check Daily Reward Status
```http
GET /api/v1/rewards/daily
Authorization: Bearer <jwt-token>
```

### Claim Daily Reward
```http
POST /api/v1/rewards/daily/claim
Authorization: Bearer <jwt-token>
```

---

## 🏪 Marketplace

### Get Marketplace Listings
```http
GET /api/v1/marketplace?page=1&per_page=10
```

### Create Listing
```http
POST /api/v1/marketplace
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "item_name": "Advanced JavaScript Guide",
  "description": "Comprehensive JS tutorial",
  "price": 500,
  "quantity": 1
}
```

### Buy from Marketplace
```http
POST /api/v1/marketplace/:id/buy
Authorization: Bearer <jwt-token>
Content-Type: application/json

{
  "quantity": 1
}
```

---

## 📚 Tutorial System

### Get Available Tutorials
```http
GET /api/v1/tutorials
Authorization: Bearer <jwt-token>
```

### Start Tutorial
```http
POST /api/v1/tutorials/:id/start
Authorization: Bearer <jwt-token>
```

### Complete Tutorial
```http
POST /api/v1/tutorials/:id/complete
Authorization: Bearer <jwt-token>
```

---

## 📊 Statistics

### Get Personal Statistics
```http
GET /api/v1/stats/me
Authorization: Bearer <jwt-token>
```

---

## 🛡️ Admin Endpoints

### Get All Users
```http
GET /api/v1/admin/users?page=1&per_page=10
Authorization: Bearer <admin-jwt-token>
```

### Ban User
```http
PUT /api/v1/admin/users/:id/ban
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json

{
  "reason": "Violation of terms",
  "duration": 7
}
```

### Change User Role
```http
PUT /api/v1/admin/users/:id/role
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json

{
  "role": "admin"
}
```

### Get System Statistics
```http
GET /api/v1/admin/stats
Authorization: Bearer <admin-jwt-token>
```

### Get Audit Logs
```http
GET /api/v1/admin/logs?page=1&per_page=50
Authorization: Bearer <admin-jwt-token>
```

### Create Quest (Admin)
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

---

## 🗄 Database Schema

### Core Models

- **User**: Player accounts with authentication and game stats
- **Friendship**: Friend relationships and status
- **OnlineUser**: Track user online status
- **Inventory**: Player items (tools, code snippets, resources)
- **Quest**: Available quests with requirements and rewards
- **UserQuestProgress**: Player progress on specific quests
- **NPC**: Non-player characters with dialogue and locations
- **NPCRelationship**: Player relationships with NPCs
- **DailyTask**: Daily challenges for players
- **UserDailyTaskProgress**: Player completion of daily tasks
- **Achievement**: Unlockable achievements with conditions
- **UserAchievement**: Player's unlocked achievements
- **Badge**: Visual rewards and recognition
- **UserBadge**: Player's earned badges
- **Skill**: Available skills and upgrades
- **UserSkill**: Player's skill progression
- **StoryProgress**: Player's story chapter progress
- **CodeBattle**: Coding challenges and battles
- **MiniGame**: Mini games and coding challenges
- **MiniGameSession**: Player game sessions
- **ShopItem**: Items available for purchase
- **UserPurchase**: Purchase history
- **CraftingRecipe**: Item crafting recipes
- **CraftingSession**: Active crafting sessions
- **DailyReward**: Daily login rewards
- **UserDailyReward**: Player's claimed rewards
- **LoginStreak**: Player login streaks
- **Event**: Time-based events and festivals
- **EventParticipant**: Event participation
- **Guild**: Player communities
- **GuildMember**: Guild membership
- **GuildInvitation**: Guild invitations
- **MarketplaceListing**: Player-to-player item sales
- **MarketplaceTransaction**: Marketplace transactions
- **Tutorial**: Learning content
- **UserTutorialProgress**: Tutorial completion
- **Notification**: System notifications
- **UserStatistics**: Player analytics
- **DailyStatistics**: Daily player metrics

## 🔧 Middleware

- **Authentication**: JWT token validation
- **CORS**: Cross-origin request handling
- **Rate Limiting**: Request rate limiting per IP
- **Logging**: Request/response logging
- **Error Handling**: Centralized error handling
- **Recovery**: Panic recovery

## 📈 Performance Tips

1. **Database Indexing**: Indexes on frequently queried columns
2. **Connection Pooling**: Optimized connection pool settings
3. **WebSocket Scaling**: Horizontal scaling with Redis pub/sub
4. **Caching**: Redis for session management and caching
5. **Pagination**: Always use pagination for list endpoints
6. **Rate Limiting**: Appropriate rate limits per endpoint

## 🧪 Testing

Run tests with:
```bash
go test ./...
```

## 🚀 Deployment

### Docker Deployment
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### Environment Setup
```bash
# Production environment variables
export DB_HOST=your-db-host
export DB_PASSWORD=your-secure-password
export JWT_SECRET=your-super-secure-jwt-secret
```

## 📝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🤝 Support

For support, email support@codevalley.com or join our Discord server.

---

**Happy Coding in Code Valley! 🏘️💻**