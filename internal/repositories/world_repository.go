package repositories

import (
	"code-valley-api/internal/database"
	"code-valley-api/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorldRepository struct {
	db *gorm.DB
}

func NewWorldRepository() *WorldRepository {
	return &WorldRepository{
		db: database.GetDB(),
	}
}

// Map operations
func (r *WorldRepository) GetMapByName(name string) (*models.Map, error) {
	var mapData models.Map
	err := r.db.Where("name = ? AND is_active = ?", name, true).First(&mapData).Error
	return &mapData, err
}

func (r *WorldRepository) GetMapByID(id uuid.UUID) (*models.Map, error) {
	var mapData models.Map
	err := r.db.First(&mapData, "id = ?", id).Error
	return &mapData, err
}

// Player position operations
func (r *WorldRepository) GetPlayerPosition(userID uuid.UUID) (*models.PlayerPosition, error) {
	var position models.PlayerPosition
	err := r.db.Preload("Map").Where("user_id = ?", userID).First(&position).Error
	return &position, err
}

func (r *WorldRepository) CreatePlayerPosition(position *models.PlayerPosition) error {
	return r.db.Create(position).Error
}

func (r *WorldRepository) UpdatePlayerPosition(position *models.PlayerPosition) error {
	return r.db.Save(position).Error
}

func (r *WorldRepository) GetPlayersInMap(mapID uuid.UUID) ([]models.PlayerPosition, error) {
	var positions []models.PlayerPosition
	err := r.db.Preload("User").Where("map_id = ?", mapID).Find(&positions).Error
	return positions, err
}

// World object operations
func (r *WorldRepository) GetWorldObjects(mapID uuid.UUID) ([]models.WorldObject, error) {
	var objects []models.WorldObject
	err := r.db.Where("map_id = ? AND is_active = ?", mapID, true).Find(&objects).Error
	return objects, err
}

func (r *WorldRepository) GetWorldObjectsAt(mapID uuid.UUID, posX, posY int) ([]models.WorldObject, error) {
	var objects []models.WorldObject
	err := r.db.Where("map_id = ? AND pos_x = ? AND pos_y = ? AND is_active = ?", 
		mapID, posX, posY, true).Find(&objects).Error
	return objects, err
}

func (r *WorldRepository) UpdateWorldObject(obj *models.WorldObject) error {
	return r.db.Save(obj).Error
}

// NPC position operations
func (r *WorldRepository) GetNPCPositions(mapID uuid.UUID) ([]models.NPCPosition, error) {
	var positions []models.NPCPosition
	err := r.db.Preload("NPC").Where("map_id = ?", mapID).Find(&positions).Error
	return positions, err
}

func (r *WorldRepository) UpdateNPCPosition(position *models.NPCPosition) error {
	return r.db.Save(position).Error
}

func (r *WorldRepository) GetNPCSchedules(npcID uuid.UUID) ([]models.NPCSchedule, error) {
	var schedules []models.NPCSchedule
	err := r.db.Where("npc_id = ?", npcID).Find(&schedules).Error
	return schedules, err
}

// Game clock operations
func (r *WorldRepository) GetGameClock() (*models.GameClock, error) {
	var clock models.GameClock
	err := r.db.First(&clock).Error
	if err == gorm.ErrRecordNotFound {
		// Create default game clock
		clock = models.GameClock{
			GameYear:   1,
			GameSeason: "spring",
			GameDay:    1,
			GameHour:   6,
			GameMinute: 0,
			TimeScale:  1.0,
		}
		r.db.Create(&clock)
	}
	return &clock, err
}

func (r *WorldRepository) UpdateGameClock(clock *models.GameClock) error {
	return r.db.Save(clock).Error
}

// Code farm operations
func (r *WorldRepository) GetUserCodeFarms(userID uuid.UUID) ([]models.CodeFarm, error) {
	var farms []models.CodeFarm
	err := r.db.Where("user_id = ?", userID).Find(&farms).Error
	return farms, err
}

func (r *WorldRepository) GetCodeFarm(farmID uuid.UUID) (*models.CodeFarm, error) {
	var farm models.CodeFarm
	err := r.db.First(&farm, "id = ?", farmID).Error
	return &farm, err
}

func (r *WorldRepository) GetCodeFarmAt(userID uuid.UUID, plotX, plotY int) (*models.CodeFarm, error) {
	var farm models.CodeFarm
	err := r.db.Where("user_id = ? AND plot_x = ? AND plot_y = ?", userID, plotX, plotY).First(&farm).Error
	return &farm, err
}

func (r *WorldRepository) CreateCodeFarm(farm *models.CodeFarm) error {
	return r.db.Create(farm).Error
}

func (r *WorldRepository) UpdateCodeFarm(farm *models.CodeFarm) error {
	return r.db.Save(farm).Error
}

func (r *WorldRepository) DeleteCodeFarm(farmID uuid.UUID) error {
	return r.db.Delete(&models.CodeFarm{}, "id = ?", farmID).Error
}