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

func (s *UserServiceImpl) CreateUser(passportNumber string) (*dto.UserResponse, error) {
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

	return &dto.UserResponse{
		ID:             user.ID,
		PassportNumber: user.PassportNumber,
		Surname:        user.Surname,
		Name:           user.Name,
		Patronymic:     user.Patronymic,
		Address:        user.Address,
	}, nil
}

func (s *UserServiceImpl) GetUserById(id uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:             user.ID,
		PassportNumber: user.PassportNumber,
		Surname:        user.Surname,
		Name:           user.Name,
		Patronymic:     user.Patronymic,
		Address:        user.Address,
	}, nil
}

func (s *UserServiceImpl) GetAllUsers() ([]dto.UserResponse, error) {
	users, err := s.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:             user.ID,
			PassportNumber: user.PassportNumber,
			Surname:        user.Surname,
			Name:           user.Name,
			Patronymic:     user.Patronymic,
			Address:        user.Address,
		})
	}

	return userResponses, nil
}

func (s *UserServiceImpl) UpdateUser(userId uint, userUpdateRequest dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetById(userId)
	if err != nil {
		return nil, err
	}

	user.Name = userUpdateRequest.Name
	user.Surname = userUpdateRequest.Surname
	user.Patronymic = userUpdateRequest.Patronymic
	user.Address = userUpdateRequest.Address

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:             user.ID,
		PassportNumber: user.PassportNumber,
		Surname:        user.Surname,
		Name:           user.Name,
		Patronymic:     user.Patronymic,
		Address:        user.Address,
	}, nil
}

func (s *UserServiceImpl) DeleteUser(id uint) error {
	return s.userRepo.Delete(id)
}

func (s *UserServiceImpl) GetUsersWithFiltersAndPagination(filters map[string]interface{}, page int, pageSize int) ([]dto.UserResponse, error) {
	users, err := s.userRepo.GetAllWithFiltersAndPagination(filters, page, pageSize)
	if err != nil {
		return nil, err
	}

	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponse := dto.UserResponse{
			ID:             user.ID,
			PassportNumber: user.PassportNumber,
			Surname:        user.Surname,
			Name:           user.Name,
			Patronymic:     user.Patronymic,
			Address:        user.Address,
		}
		userResponses = append(userResponses, userResponse)
	}

	return userResponses, nil
}
