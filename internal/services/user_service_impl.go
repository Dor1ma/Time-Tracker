package services

import (
	"encoding/json"
	"fmt"
	"github.com/Dor1ma/Time-Tracker/internal/dto"
	"github.com/Dor1ma/Time-Tracker/internal/models"
	"github.com/Dor1ma/Time-Tracker/internal/repositories"
	"net/http"
	"strconv"
)

type UserServiceImpl struct {
	userRepo       repositories.UserRepository
	externalAPIURL string
}

func NewUserServiceImpl(userRepo repositories.UserRepository, externalAPIURL string) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo:       userRepo,
		externalAPIURL: externalAPIURL,
	}
}

func (s *UserServiceImpl) CreateUser(passportNumber string) (*models.User, error) {
	passportSerie, err := strconv.Atoi(passportNumber[:4])
	if err != nil {
		return nil, err
	}
	passportNum, err := strconv.Atoi(passportNumber[5:])
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/info?passportSerie=%d&passportNumber=%d", s.externalAPIURL, passportSerie, passportNum)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResponse dto.ExternalAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	user := &models.User{
		PassportNumber: passportNumber,
		Surname:        apiResponse.Surname,
		Name:           apiResponse.Name,
		Patronymic:     apiResponse.Patronymic,
		Address:        apiResponse.Address,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserServiceImpl) GetUserById(id uint) (*models.User, error) {
	return s.userRepo.GetById(id)
}

func (s *UserServiceImpl) GetAllUsers() (*[]models.User, error) {
	return s.userRepo.GetAll()
}

func (s *UserServiceImpl) GetUsersWithFilters(filters map[string]interface{}) (*[]models.User, error) {
	return s.userRepo.GetAllWithFilters(filters)
}

func (s *UserServiceImpl) GetUsersWithPagination(page int, pageSize int) (*[]models.User, error) {
	return s.userRepo.GetWithPagination(page, pageSize)
}

func (s *UserServiceImpl) UpdateUser(user *models.User) error {
	return s.userRepo.Update(user)
}

func (s *UserServiceImpl) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}
