package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	clients    map[*Client]bool
	userClients map[uuid.UUID]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
	
	// Channels for handling game events
	playerMoveChannel chan PlayerMoveEvent
	playerInteractChannel chan PlayerInteractEvent
}

type PlayerMoveEvent struct {
	UserID    uuid.UUID
	PosX      int
	PosY      int
	Direction string
}

type PlayerInteractEvent struct {
	UserID  uuid.UUID
	TargetX int
	TargetY int
}

type MapClients struct {
	clients map[uuid.UUID]map[*Client]bool
	mutex   sync.RWMutex
}

var mapClients = &MapClients{
	clients: make(map[uuid.UUID]map[*Client]bool),
}

type Message struct {
	Type    string      `json:"type"`
	UserID  uuid.UUID   `json:"user_id,omitempty"`
	Data    interface{} `json:"data"`
	Target  string      `json:"target,omitempty"` // "all", "user", "friends"
}

type PlayerMoveMessage struct {
	Type      string `json:"type"`
	PosX      int    `json:"pos_x"`
	PosY      int    `json:"pos_y"`
	Direction string `json:"direction"`
}

type PlayerInteractMessage struct {
	Type    string    `json:"type"`
	TargetX int       `json:"target_x"`
	TargetY int       `json:"target_y"`
	TargetID *uuid.UUID `json:"target_id,omitempty"`
}

func NewHub() *Hub {
	return &Hub{
		clients:               make(map[*Client]bool),
		userClients:           make(map[uuid.UUID]*Client),
		broadcast:             make(chan []byte, 256),
		register:              make(chan *Client),
		unregister:            make(chan *Client),
		playerMoveChannel:     make(chan PlayerMoveEvent, 256),
		playerInteractChannel: make(chan PlayerInteractEvent, 256),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.userClients[client.UserID] = client
			h.mutex.Unlock()
			
			log.Printf("Client connected: %s", client.UserID)
			
			// Notify friends that user is online
			h.notifyFriendsOnlineStatus(client.UserID, true)

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.userClients, client.UserID)
				close(client.send)
			}
			h.mutex.Unlock()
			
			log.Printf("Client disconnected: %s", client.UserID)
			
			// Notify friends that user is offline
			h.notifyFriendsOnlineStatus(client.UserID, false)

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
					delete(h.userClients, client.UserID)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

func (h *Hub) SendToUser(userID uuid.UUID, message Message) {
	h.mutex.RLock()
	client, exists := h.userClients[userID]
	h.mutex.RUnlock()
	
	if exists {
		data, err := json.Marshal(message)
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			return
		}
		
		select {
		case client.send <- data:
		default:
			h.mutex.Lock()
			close(client.send)
			delete(h.clients, client)
			delete(h.userClients, userID)
			h.mutex.Unlock()
		}
	}
}

func (h *Hub) SendToAll(message Message) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}
	
	select {
	case h.broadcast <- data:
	default:
		log.Println("Broadcast channel is full")
	}
}

func (h *Hub) GetOnlineUsers() []uuid.UUID {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	users := make([]uuid.UUID, 0, len(h.userClients))
	for userID := range h.userClients {
		users = append(users, userID)
	}
	return users
}

func (h *Hub) IsUserOnline(userID uuid.UUID) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	_, exists := h.userClients[userID]
	return exists
}

func (h *Hub) AddClientToMap(client *Client, mapID uuid.UUID) {
	mapClients.mutex.Lock()
	defer mapClients.mutex.Unlock()
	
	if mapClients.clients[mapID] == nil {
		mapClients.clients[mapID] = make(map[*Client]bool)
	}
	mapClients.clients[mapID][client] = true
}

func (h *Hub) RemoveClientFromMap(client *Client, mapID uuid.UUID) {
	mapClients.mutex.Lock()
	defer mapClients.mutex.Unlock()
	
	if mapClients.clients[mapID] != nil {
		delete(mapClients.clients[mapID], client)
		if len(mapClients.clients[mapID]) == 0 {
			delete(mapClients.clients, mapID)
		}
	}
}

func BroadcastToMap(mapID uuid.UUID, message Message) {
	mapClients.mutex.RLock()
	clients := mapClients.clients[mapID]
	mapClients.mutex.RUnlock()
	
	if clients == nil {
		return
	}
	
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}
	
	for client := range clients {
		select {
		case client.send <- data:
		default:
			// Client channel is full, remove it
			close(client.send)
			delete(clients, client)
		}
	}
}

func (h *Hub) notifyFriendsOnlineStatus(userID uuid.UUID, isOnline bool) {
	// This would typically query the database for friends
	// For now, we'll just broadcast to all users
	message := Message{
		Type: "user_status",
		Data: map[string]interface{}{
			"user_id":   userID,
			"is_online": isOnline,
		},
	}
	
	h.SendToAll(message)
}

func (h *Hub) HandlePlayerMove(userID uuid.UUID, posX, posY int, direction string) {
	select {
	case h.playerMoveChannel <- PlayerMoveEvent{
		UserID:    userID,
		PosX:      posX,
		PosY:      posY,
		Direction: direction,
	}:
	default:
		log.Println("Player move channel is full")
	}
}

func (h *Hub) HandlePlayerInteract(userID uuid.UUID, targetX, targetY int) {
	select {
	case h.playerInteractChannel <- PlayerInteractEvent{
		UserID:  userID,
		TargetX: targetX,
		TargetY: targetY,
	}:
	default:
		log.Println("Player interact channel is full")
	}
}

func (h *Hub) GetPlayerMoveChannel() <-chan PlayerMoveEvent {
	return h.playerMoveChannel
}

func (h *Hub) GetPlayerInteractChannel() <-chan PlayerInteractEvent {
	return h.playerInteractChannel
}