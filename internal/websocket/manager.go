package websocket

import (
	"log"

	"code-valley-api/internal/config"
	"code-valley-api/internal/utils"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var GlobalHub *Hub

func InitializeWebSocket() {
	GlobalHub = NewHub()
	go GlobalHub.Run()
	log.Println("WebSocket hub initialized")
}

func WebSocketUpgrade(cfg *config.Config) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		// Get user from query params or headers
		token := c.Query("token")
		if token == "" {
			c.Close()
			return
		}

		claims, err := utils.ValidateJWT(token, cfg.JWT.Secret)
		if err != nil {
			log.Printf("WebSocket auth failed: %v", err)
			c.Close()
			return
		}

		client := NewClient(GlobalHub, c, claims.UserID)
		GlobalHub.register <- client

		go client.WritePump()
		client.ReadPump()
	})
}

// Utility functions for sending real-time updates
func NotifyQuestUpdate(userID uuid.UUID, questData interface{}) {
	if GlobalHub != nil {
		message := Message{
			Type:   "quest_update",
			UserID: userID,
			Data:   questData,
		}
		GlobalHub.SendToUser(userID, message)
	}
}

func NotifyFriendRequest(userID uuid.UUID, friendData interface{}) {
	if GlobalHub != nil {
		message := Message{
			Type:   "friend_request",
			UserID: userID,
			Data:   friendData,
		}
		GlobalHub.SendToUser(userID, message)
	}
}

func NotifyAchievementUnlocked(userID uuid.UUID, achievementData interface{}) {
	if GlobalHub != nil {
		message := Message{
			Type:   "achievement_unlocked",
			UserID: userID,
			Data:   achievementData,
		}
		GlobalHub.SendToUser(userID, message)
	}
}

func NotifyLevelUp(userID uuid.UUID, levelData interface{}) {
	if GlobalHub != nil {
		message := Message{
			Type:   "level_up",
			UserID: userID,
			Data:   levelData,
		}
		GlobalHub.SendToUser(userID, message)
	}
}

func BroadcastEvent(eventData interface{}) {
	if GlobalHub != nil {
		message := Message{
			Type: "event_broadcast",
			Data: eventData,
		}
		GlobalHub.SendToAll(message)
	}
}