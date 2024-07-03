package repositories

import "github.com/Dor1ma/Time-Tracker/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	GetById(id uint) (*models.User, error)
	GetAll() ([]models.User, error)
	GetAllWithFilters(filters map[string]interface{}) ([]models.User, error)
	GetWithPagination(page int, pageSize int) ([]models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
}
