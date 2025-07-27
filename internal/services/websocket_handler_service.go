package services

import (
	"log"

	"code-valley-api/internal/websocket"
)

type WebSocketHandlerService struct {
	worldService *WorldService
}

func NewWebSocketHandlerService() *WebSocketHandlerService {
	return &WebSocketHandlerService{
		worldService: NewWorldService(),
	}
}

func (s *WebSocketHandlerService) Start() {
	go s.handlePlayerMoves()
	go s.handlePlayerInteractions()
	log.Println("WebSocket handler service started")
}

func (s *WebSocketHandlerService) handlePlayerMoves() {
	for moveEvent := range websocket.GlobalHub.GetPlayerMoveChannel() {
		err := s.worldService.MovePlayer(moveEvent.UserID, moveEvent.PosX, moveEvent.PosY, moveEvent.Direction)
		if err != nil {
			// Send error back to client
			websocket.GlobalHub.SendToUser(moveEvent.UserID, websocket.Message{
				Type: "movement_error",
				Data: map[string]interface{}{
					"error": err.Error(),
				},
			})
		}
	}
}

func (s *WebSocketHandlerService) handlePlayerInteractions() {
	for interactEvent := range websocket.GlobalHub.GetPlayerInteractChannel() {
		result, err := s.worldService.InteractWithObject(interactEvent.UserID, interactEvent.TargetX, interactEvent.TargetY)
		
		response := websocket.Message{
			Type: "interaction_result",
			Data: map[string]interface{}{
				"success": err == nil,
				"result":  result,
			},
		}
		
		if err != nil {
			response.Data.(map[string]interface{})["error"] = err.Error()
		}
		
		websocket.GlobalHub.SendToUser(interactEvent.UserID, response)
	}
}