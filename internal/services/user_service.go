package services

import "github.com/Dor1ma/Time-Tracker/internal/models"

type UserService interface {
	CreateUser(passportNumber string) (*models.User, error)
	GetUserById(userId uint) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	GetUsersWithFilters(filters map[string]interface{}) ([]models.User, error)
	GetUsersWithPagination(page int, pageSize int) ([]models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
}
