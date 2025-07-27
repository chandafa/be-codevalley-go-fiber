package services

import (
	"errors"
	"time"

	"code-valley-api/internal/models"
	"code-valley-api/internal/repositories"
	"code-valley-api/internal/utils"
	"code-valley-api/internal/websocket"

	"github.com/google/uuid"
)

type WorldService struct {
	worldRepo     *repositories.WorldRepository
	userRepo      *repositories.UserRepository
	inventoryRepo *repositories.InventoryRepository
}

func NewWorldService() *WorldService {
	return &WorldService{
		worldRepo:     repositories.NewWorldRepository(),
		userRepo:      repositories.NewUserRepository(),
		inventoryRepo: repositories.NewInventoryRepository(),
	}
}

type MapStateResponse struct {
	Map             *models.Map             `json:"map"`
	PlayerPositions []models.PlayerPosition `json:"player_positions"`
	WorldObjects    []models.WorldObject    `json:"world_objects"`
	NPCPositions    []models.NPCPosition    `json:"npc_positions"`
	GameTime        *models.GameClock       `json:"game_time"`
}

func (s *WorldService) GetMapState(mapName string, userID uuid.UUID) (*MapStateResponse, error) {
	mapData, err := s.worldRepo.GetMapByName(mapName)
	if err != nil {
		return nil, errors.New("map not found")
	}

	playerPositions, err := s.worldRepo.GetPlayersInMap(mapData.ID)
	if err != nil {
		return nil, err
	}

	worldObjects, err := s.worldRepo.GetWorldObjects(mapData.ID)
	if err != nil {
		return nil, err
	}

	npcPositions, err := s.worldRepo.GetNPCPositions(mapData.ID)
	if err != nil {
		return nil, err
	}

	gameTime, err := s.worldRepo.GetGameClock()
	if err != nil {
		return nil, err
	}

	// Ensure player has a position in this map
	_, err = s.worldRepo.GetPlayerPosition(userID)
	if err != nil {
		// Create default position for new player
		defaultPosition := &models.PlayerPosition{
			UserID:    userID,
			MapID:     mapData.ID,
			PosX:      mapData.Width / 2,
			PosY:      mapData.Height / 2,
			Direction: "down",
			LastMoved: time.Now(),
		}
		s.worldRepo.CreatePlayerPosition(defaultPosition)
	}

	return &MapStateResponse{
		Map:             mapData,
		PlayerPositions: playerPositions,
		WorldObjects:    worldObjects,
		NPCPositions:    npcPositions,
		GameTime:        gameTime,
	}, nil
}

func (s *WorldService) GetPlayerPosition(userID uuid.UUID) (*models.PlayerPosition, error) {
	return s.worldRepo.GetPlayerPosition(userID)
}

type TeleportRequest struct {
	MapName string `json:"map_name" validate:"required"`
	PosX    int    `json:"pos_x" validate:"min=0"`
	PosY    int    `json:"pos_y" validate:"min=0"`
}

func (s *WorldService) TeleportPlayer(userID uuid.UUID, req TeleportRequest) (*models.PlayerPosition, error) {
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	mapData, err := s.worldRepo.GetMapByName(req.MapName)
	if err != nil {
		return nil, errors.New("map not found")
	}

	if req.PosX >= mapData.Width || req.PosY >= mapData.Height {
		return nil, errors.New("position out of bounds")
	}

	position, err := s.worldRepo.GetPlayerPosition(userID)
	if err != nil {
		// Create new position
		position = &models.PlayerPosition{
			UserID:    userID,
			MapID:     mapData.ID,
			PosX:      req.PosX,
			PosY:      req.PosY,
			Direction: "down",
			LastMoved: time.Now(),
		}
		err = s.worldRepo.CreatePlayerPosition(position)
	} else {
		// Update existing position
		position.MapID = mapData.ID
		position.PosX = req.PosX
		position.PosY = req.PosY
		position.LastMoved = time.Now()
		err = s.worldRepo.UpdatePlayerPosition(position)
	}

	if err != nil {
		return nil, err
	}

	// Broadcast position update
	websocket.BroadcastToMap(mapData.ID, websocket.Message{
		Type: "player_position_update",
		Data: map[string]interface{}{
			"user_id":   userID,
			"map_id":    mapData.ID,
			"pos_x":     req.PosX,
			"pos_y":     req.PosY,
			"direction": position.Direction,
		},
	})

	return position, nil
}

func (s *WorldService) MovePlayer(userID uuid.UUID, posX, posY int, direction string) error {
	position, err := s.worldRepo.GetPlayerPosition(userID)
	if err != nil {
		return errors.New("player position not found")
	}

	// Validate movement (basic bounds checking)
	mapData, err := s.worldRepo.GetMapByID(position.MapID)
	if err != nil {
		return errors.New("map not found")
	}

	if posX < 0 || posX >= mapData.Width || posY < 0 || posY >= mapData.Height {
		return errors.New("position out of bounds")
	}

	// Check for collisions with world objects
	objects, err := s.worldRepo.GetWorldObjectsAt(position.MapID, posX, posY)
	if err == nil && len(objects) > 0 {
		for _, obj := range objects {
			if obj.ObjectType == models.ObjectTypeRock || obj.ObjectType == models.ObjectTypeTree {
				return errors.New("position blocked by object")
			}
		}
	}

	// Update position
	position.PosX = posX
	position.PosY = posY
	position.Direction = direction
	position.LastMoved = time.Now()

	if err := s.worldRepo.UpdatePlayerPosition(position); err != nil {
		return err
	}

	// Broadcast movement
	websocket.BroadcastToMap(position.MapID, websocket.Message{
		Type: "player_position_update",
		Data: map[string]interface{}{
			"user_id":   userID,
			"map_id":    position.MapID,
			"pos_x":     posX,
			"pos_y":     posY,
			"direction": direction,
		},
	})

	return nil
}

func (s *WorldService) InteractWithObject(userID uuid.UUID, targetX, targetY int) (map[string]interface{}, error) {
	position, err := s.worldRepo.GetPlayerPosition(userID)
	if err != nil {
		return nil, errors.New("player position not found")
	}

	// Check if target is within interaction range (adjacent tiles)
	if abs(position.PosX-targetX) > 1 || abs(position.PosY-targetY) > 1 {
		return nil, errors.New("target too far away")
	}

	// Get objects at target position
	objects, err := s.worldRepo.GetWorldObjectsAt(position.MapID, targetX, targetY)
	if err != nil || len(objects) == 0 {
		return nil, errors.New("no interactable object found")
	}

	obj := objects[0]
	result := make(map[string]interface{})

	switch obj.ObjectType {
	case models.ObjectTypeTree:
		result = s.chopTree(userID, &obj)
	case models.ObjectTypeRock:
		result = s.mineRock(userID, &obj)
	case models.ObjectTypeChest:
		result = s.openChest(userID, &obj)
	case models.ObjectTypeServer:
		result = s.accessServer(userID, &obj)
	default:
		return nil, errors.New("object not interactable")
	}

	// Update object state
	s.worldRepo.UpdateWorldObject(&obj)

	// Broadcast object update
	websocket.BroadcastToMap(position.MapID, websocket.Message{
		Type: "world_object_update",
		Data: map[string]interface{}{
			"object_id": obj.ID,
			"pos_x":     obj.PosX,
			"pos_y":     obj.PosY,
			"state":     obj.State,
		},
	})

	return result, nil
}

func (s *WorldService) GetGameTime() (*models.GameClock, error) {
	return s.worldRepo.GetGameClock()
}

// Code Farming System
type PlantCodeRequest struct {
	PlotX    int    `json:"plot_x" validate:"min=0"`
	PlotY    int    `json:"plot_y" validate:"min=0"`
	CodeType string `json:"code_type" validate:"required"`
}

func (s *WorldService) GetUserCodeFarms(userID uuid.UUID) ([]models.CodeFarm, error) {
	return s.worldRepo.GetUserCodeFarms(userID)
}

func (s *WorldService) PlantCode(userID uuid.UUID, req PlantCodeRequest) (*models.CodeFarm, error) {
	if err := utils.ValidateStruct(req); err != nil {
		return nil, err
	}

	// Check if plot is already occupied
	existing, err := s.worldRepo.GetCodeFarmAt(userID, req.PlotX, req.PlotY)
	if err == nil && existing != nil {
		return nil, errors.New("plot already occupied")
	}

	now := time.Now()
	harvestTime := now.Add(24 * time.Hour) // 1 day to grow

	farm := &models.CodeFarm{
		UserID:      userID,
		PlotX:       req.PlotX,
		PlotY:       req.PlotY,
		CodeType:    req.CodeType,
		PlantedAt:   &now,
		LastWatered: &now,
		HarvestAt:   &harvestTime,
		GrowthStage: 1,
		Quality:     "normal",
	}

	if err := s.worldRepo.CreateCodeFarm(farm); err != nil {
		return nil, err
	}

	return farm, nil
}

func (s *WorldService) WaterCode(userID uuid.UUID, farmID uuid.UUID) (*models.CodeFarm, error) {
	farm, err := s.worldRepo.GetCodeFarm(farmID)
	if err != nil {
		return nil, errors.New("code farm not found")
	}

	if farm.UserID != userID {
		return nil, errors.New("not your code farm")
	}

	now := time.Now()
	farm.LastWatered = &now

	// Improve growth if watered regularly
	if farm.LastWatered != nil && now.Sub(*farm.LastWatered) < 12*time.Hour {
		if farm.GrowthStage < 4 {
			farm.GrowthStage++
		}
	}

	if err := s.worldRepo.UpdateCodeFarm(farm); err != nil {
		return nil, err
	}

	return farm, nil
}

func (s *WorldService) HarvestCode(userID uuid.UUID, farmID uuid.UUID) (map[string]interface{}, error) {
	farm, err := s.worldRepo.GetCodeFarm(farmID)
	if err != nil {
		return nil, errors.New("code farm not found")
	}

	if farm.UserID != userID {
		return nil, errors.New("not your code farm")
	}

	if farm.HarvestAt == nil || time.Now().Before(*farm.HarvestAt) {
		return nil, errors.New("code not ready for harvest")
	}

	// Calculate harvest rewards based on quality and growth stage
	baseReward := 50
	qualityMultiplier := 1.0
	switch farm.Quality {
	case "silver":
		qualityMultiplier = 1.25
	case "gold":
		qualityMultiplier = 1.5
	case "iridium":
		qualityMultiplier = 2.0
	}

	coins := int(float64(baseReward) * qualityMultiplier * float64(farm.GrowthStage))
	exp := coins / 2

	// Add rewards to user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	user.Coins += coins
	user.EXP += exp
	s.userRepo.Update(user)

	// Add harvested item to inventory
	item := &models.Inventory{
		UserID:   userID,
		ItemName: farm.CodeType + " Library",
		Quantity: 1,
		ItemType: models.ItemTypeCode,
	}
	s.inventoryRepo.AddItem(item)

	// Remove the farm
	s.worldRepo.DeleteCodeFarm(farmID)

	return map[string]interface{}{
		"coins_earned": coins,
		"exp_earned":   exp,
		"item_name":    item.ItemName,
		"quality":      farm.Quality,
	}, nil
}

// Helper functions for object interactions
func (s *WorldService) chopTree(userID uuid.UUID, obj *models.WorldObject) map[string]interface{} {
	hp := 3
	if obj.State["hp"] != nil {
		hp = int(obj.State["hp"].(float64))
	}

	hp--
	obj.State["hp"] = hp

	if hp <= 0 {
		obj.IsActive = false
		// Add wood to inventory
		item := &models.Inventory{
			UserID:   userID,
			ItemName: "Code Snippet",
			Quantity: 2,
			ItemType: models.ItemTypeResource,
		}
		s.inventoryRepo.AddItem(item)

		return map[string]interface{}{
			"action":       "tree_chopped",
			"items_gained": []string{"Code Snippet x2"},
		}
	}

	return map[string]interface{}{
		"action":        "tree_damaged",
		"remaining_hp":  hp,
	}
}

func (s *WorldService) mineRock(userID uuid.UUID, obj *models.WorldObject) map[string]interface{} {
	hp := 2
	if obj.State["hp"] != nil {
		hp = int(obj.State["hp"].(float64))
	}

	hp--
	obj.State["hp"] = hp

	if hp <= 0 {
		obj.IsActive = false
		// Add ore to inventory
		item := &models.Inventory{
			UserID:   userID,
			ItemName: "Raw Data",
			Quantity: 1,
			ItemType: models.ItemTypeResource,
		}
		s.inventoryRepo.AddItem(item)

		return map[string]interface{}{
			"action":       "rock_mined",
			"items_gained": []string{"Raw Data x1"},
		}
	}

	return map[string]interface{}{
		"action":       "rock_damaged",
		"remaining_hp": hp,
	}
}

func (s *WorldService) openChest(userID uuid.UUID, obj *models.WorldObject) map[string]interface{} {
	if obj.State["is_looted"] == true {
		return map[string]interface{}{
			"action":  "chest_empty",
			"message": "This chest has already been looted",
		}
	}

	obj.State["is_looted"] = true

	// Add random item to inventory
	items := []string{"Debug Tool", "Refactor Kit", "Unit Test Template"}
	randomItem := items[time.Now().Unix()%int64(len(items))]

	item := &models.Inventory{
		UserID:   userID,
		ItemName: randomItem,
		Quantity: 1,
		ItemType: models.ItemTypeTool,
	}
	s.inventoryRepo.AddItem(item)

	return map[string]interface{}{
		"action":       "chest_opened",
		"items_gained": []string{randomItem + " x1"},
	}
}

func (s *WorldService) accessServer(userID uuid.UUID, obj *models.WorldObject) map[string]interface{} {
	return map[string]interface{}{
		"action":  "server_accessed",
		"message": "Connected to server. You can now deploy your code!",
		"options": []string{"Deploy", "Monitor", "Scale"},
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}