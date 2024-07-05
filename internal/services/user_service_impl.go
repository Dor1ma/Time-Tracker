package services

import (
	"encoding/json"
	"fmt"
	"github.com/Dor1ma/Time-Tracker/internal/dto"
	"github.com/Dor1ma/Time-Tracker/internal/models"
	"github.com/Dor1ma/Time-Tracker/internal/repositories"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type UserServiceImpl struct {
	userRepo       repositories.UserRepository
	externalAPIURL string
	logger         *logrus.Logger
}

func NewUserServiceImpl(userRepo repositories.UserRepository, externalAPIURL string, logger *logrus.Logger) *UserServiceImpl {
	return &UserServiceImpl{
		userRepo:       userRepo,
		externalAPIURL: externalAPIURL,
		logger:         logger,
	}
}

func (s *UserServiceImpl) CreateUser(passportNumber string) (*dto.UserResponse, error) {
	s.logger.Infof("): creating user with passport number: %s", passportNumber)
	passportSerie, err := strconv.Atoi(passportNumber[:4])
	if err != nil {
		s.logger.Debugf("CreateUser: invalid passport series: %v", err)
		return nil, err
	}
	passportNum, err := strconv.Atoi(passportNumber[5:])
	if err != nil {
		s.logger.Debugf("CreateUser: invalid passport number: %v", err)
		return nil, err
	}

	url := fmt.Sprintf("%s/info?passportSerie=%d&passportNumber=%d", s.externalAPIURL, passportSerie, passportNum)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		s.logger.Debugf("CreateUser: failed to fetch data from external API: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	var apiResponse dto.ExternalAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		s.logger.Debugf("CreateUser: failed to decode external API response: %v", err)
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
		s.logger.Debugf("CreateUser: failed to create user in database: %v", err)
		return nil, err
	}

	s.logger.Infof("CreateUser: user created with ID: %d", user.ID)
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
	s.logger.Infof("GetUserById: getting user with id: %d", id)
	user, err := s.userRepo.GetById(id)
	if err != nil {
		s.logger.Debugf("GetUserById: failed to get user in database: %v", err)
		return nil, err
	}

	s.logger.Infof("GetUserById: got user with ID: %d", user.ID)
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
	s.logger.Info("GetAllUsers: fetching all users")
	users, err := s.userRepo.GetAll()
	if err != nil {
		s.logger.Debugf("GetAllUsers: failed to fetch users: %v", err)
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
	s.logger.Infof("UpdateUser: updating user with ID: %d", userId)
	user, err := s.userRepo.GetById(userId)
	if err != nil {
		s.logger.Errorf("UpdateUser: failed to update user: %v", err)
		return nil, err
	}

	user.Name = userUpdateRequest.Name
	user.Surname = userUpdateRequest.Surname
	user.Patronymic = userUpdateRequest.Patronymic
	user.Address = userUpdateRequest.Address

	if err := s.userRepo.Update(user); err != nil {
		s.logger.Errorf("UpdateUser: failed to update user: %v", err)
		return nil, err
	}

	s.logger.Infof("UpdateUser: user updated with ID: %d", user.ID)
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
	s.logger.Infof("DeleteUser: deleting user with ID: %d", id)
	if err := s.userRepo.Delete(id); err != nil {
		s.logger.Debugf("DeleteUser: failed to delete user: %v", err)
		return err
	}
	s.logger.Infof("DeleteUser: user deleted with ID: %d", id)
	return nil
}

func (s *UserServiceImpl) GetUsersWithFiltersAndPagination(filters map[string]interface{}, page int, pageSize int) ([]dto.UserResponse, error) {
	s.logger.Infof("GetUsersWithFiltersAndPagination: fetching users with filters and pagination: filters=%d page=%d, pageSize=%d",
		len(filters), page, pageSize)
	users, err := s.userRepo.GetAllWithFiltersAndPagination(filters, page, pageSize)
	if err != nil {
		s.logger.Debugf("GetUsersWithFiltersAndPagination: failed to fetch users: %v", err)
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
