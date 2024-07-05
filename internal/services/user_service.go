package services

import (
	"github.com/Dor1ma/Time-Tracker/internal/dto"
)

type UserService interface {
	CreateUser(passportNumber string) (*dto.UserResponse, error)
	GetUserById(userId uint) (*dto.UserResponse, error)
	GetAllUsers() ([]dto.UserResponse, error)
	GetUsersWithFiltersAndPagination(filters map[string]interface{}, page int, pageSize int) ([]dto.UserResponse, error)
	UpdateUser(userId uint, userUpdateRequest dto.UpdateUserRequest) (*dto.UserResponse, error)
	DeleteUser(id uint) error
}
