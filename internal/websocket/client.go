package websocket

import (
	"encoding/json"
	"log"
	"code-valley-api/internal/services"

	// "net/http"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	UserID uuid.UUID
	worldService *services.WorldService
}

func NewClient(hub *Hub, conn *websocket.Conn, userID uuid.UUID, worldService *services.WorldService) *Client {
	return &Client{
		hub:          hub,
		conn:         conn,
		send:         make(chan []byte, 256),
		UserID:       userID,
		worldService: worldService,
	}
}

func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var msg Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle incoming messages
		c.handleMessage(msg)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client) handleMessage(msg Message) {
	switch msg.Type {
	case "ping":
		response := Message{
			Type: "pong",
			Data: map[string]interface{}{
				"timestamp": time.Now().Unix(),
			},
		}
		data, _ := json.Marshal(response)
		select {
		case c.send <- data:
		default:
			close(c.send)
		}

	case "player_move":
		// Handle player movement
		if moveData, ok := msg.Data.(map[string]interface{}); ok {
			posX := int(moveData["pos_x"].(float64))
			posY := int(moveData["pos_y"].(float64))
			direction := moveData["direction"].(string)
			
			err := c.worldService.MovePlayer(c.UserID, posX, posY, direction)
			if err != nil {
				// Send error back to client
				response := Message{
					Type: "movement_error",
					Data: map[string]interface{}{
						"error": err.Error(),
					},
				}
				data, _ := json.Marshal(response)
				select {
				case c.send <- data:
				default:
					close(c.send)
				}
			}
		}

	case "player_interact":
		// Handle player interaction
		if interactData, ok := msg.Data.(map[string]interface{}); ok {
			targetX := int(interactData["target_x"].(float64))
			targetY := int(interactData["target_y"].(float64))
			
			result, err := c.worldService.InteractWithObject(c.UserID, targetX, targetY)
			
			response := Message{
				Type: "interaction_result",
				Data: map[string]interface{}{
					"success": err == nil,
					"result":  result,
				},
			}
			
			if err != nil {
				response.Data.(map[string]interface{})["error"] = err.Error()
			}
			
			data, _ := json.Marshal(response)
			select {
			case c.send <- data:
			default:
				close(c.send)
			}
		}

	case "chat":
		// Handle chat messages
		log.Printf("Chat message from %s: %v", c.UserID, msg.Data)

	case "quest_update":
		// Handle quest updates
		log.Printf("Quest update from %s: %v", c.UserID, msg.Data)

	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}
