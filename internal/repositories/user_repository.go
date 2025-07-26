package repositories

import (
	"code-valley-api/internal/database"
	"code-valley-api/internal/models"
	"code-valley-api/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.GetDB(),
	}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) GetAll(pagination utils.PaginationParams) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Count total records
	r.db.Model(&models.User{}).Count(&total)

	// Get paginated records
	err := r.db.Offset(pagination.Offset).
		Limit(pagination.PerPage).
		Find(&users).Error

	return users, total, err
}

func (r *UserRepository) UpdateEXPAndLevel(userID uuid.UUID, exp, level int) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"exp":   exp,
			"level": level,
		}).Error
}

func (r *UserRepository) UpdateCoins(userID uuid.UUID, coins int) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("coins", coins).Error
}