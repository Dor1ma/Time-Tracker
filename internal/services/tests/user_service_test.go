package tests

import (
	"encoding/json"
	"errors"
	"github.com/Dor1ma/Time-Tracker/internal/dto"
	"github.com/Dor1ma/Time-Tracker/internal/models"
	"github.com/Dor1ma/Time-Tracker/internal/repositories"
	"github.com/Dor1ma/Time-Tracker/internal/services"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type UserServiceTestSuite struct {
	suite.Suite
	userService     *services.UserServiceImpl
	userRepoMock    *repositories.MockUserRepository
	externalAPIMock *httptest.Server
}

func (suite *UserServiceTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.userRepoMock = repositories.NewMockUserRepository(ctrl)

	logger := logrus.New()
	suite.externalAPIMock = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiResponse := dto.ExternalAPIResponse{
			Surname:    "Doe",
			Name:       "John",
			Patronymic: "Michael",
			Address:    "123 Main St",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(apiResponse)
	}))

	suite.userService = services.NewUserServiceImpl(suite.userRepoMock, suite.externalAPIMock.URL, logger)
}

func (suite *UserServiceTestSuite) TearDownTest() {
	suite.externalAPIMock.Close()
}

func (suite *UserServiceTestSuite) TestCreateUserSuccess() {
	passportNumber := "1234-567890"
	expectedUser := &models.User{
		PassportNumber: passportNumber,
		Surname:        "Doe",
		Name:           "John",
		Patronymic:     "Michael",
		Address:        "123 Main St",
	}

	suite.userRepoMock.EXPECT().
		Create(gomock.Any()).
		DoAndReturn(func(user *models.User) error {
			user.ID = 1
			return nil
		})

	userResponse, err := suite.userService.CreateUser(passportNumber)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), userResponse)
	assert.Equal(suite.T(), uint(1), userResponse.ID)
	assert.Equal(suite.T(), expectedUser.Surname, userResponse.Surname)
	assert.Equal(suite.T(), expectedUser.Name, userResponse.Name)
	assert.Equal(suite.T(), expectedUser.Patronymic, userResponse.Patronymic)
	assert.Equal(suite.T(), expectedUser.Address, userResponse.Address)
}

func (suite *UserServiceTestSuite) TestCreateUserInvalidPassportNumber() {
	passportNumber := "invalid"
	userResponse, err := suite.userService.CreateUser(passportNumber)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), userResponse)
}

func (suite *UserServiceTestSuite) TestCreateUserExternalAPIError() {
	suite.externalAPIMock.Close()

	passportNumber := "1234-567890"
	userResponse, err := suite.userService.CreateUser(passportNumber)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), userResponse)
}

func (suite *UserServiceTestSuite) TestGetUserByIdSuccess() {
	expectedUser := &models.User{
		ID:             1,
		PassportNumber: "1234-567890",
		Surname:        "Doe",
		Name:           "John",
		Patronymic:     "Michael",
		Address:        "123 Main St",
	}

	suite.userRepoMock.EXPECT().
		GetById(uint(1)).
		Return(expectedUser, nil)

	userResponse, err := suite.userService.GetUserById(1)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), userResponse)
	assert.Equal(suite.T(), expectedUser.ID, userResponse.ID)
	assert.Equal(suite.T(), expectedUser.PassportNumber, userResponse.PassportNumber)
	assert.Equal(suite.T(), expectedUser.Surname, userResponse.Surname)
	assert.Equal(suite.T(), expectedUser.Name, userResponse.Name)
	assert.Equal(suite.T(), expectedUser.Patronymic, userResponse.Patronymic)
	assert.Equal(suite.T(), expectedUser.Address, userResponse.Address)
}

func (suite *UserServiceTestSuite) TestGetUserByIdNotFound() {
	suite.userRepoMock.EXPECT().
		GetById(uint(1)).
		Return(nil, errors.New("user not found"))

	userResponse, err := suite.userService.GetUserById(1)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), userResponse)
}

func (suite *UserServiceTestSuite) TestGetAllUsersSuccess() {
	expectedUsers := []models.User{
		{
			ID:             1,
			PassportNumber: "1234-567890",
			Surname:        "Doe",
			Name:           "John",
			Patronymic:     "Michael",
			Address:        "123 Main St",
		},
		{
			ID:             2,
			PassportNumber: "5678-123456",
			Surname:        "Smith",
			Name:           "Jane",
			Patronymic:     "Ann",
			Address:        "456 Elm St",
		},
	}

	suite.userRepoMock.EXPECT().
		GetAll().
		Return(expectedUsers, nil)

	userResponses, err := suite.userService.GetAllUsers()
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), userResponses)
	assert.Len(suite.T(), userResponses, 2)
	assert.Equal(suite.T(), expectedUsers[0].ID, userResponses[0].ID)
	assert.Equal(suite.T(), expectedUsers[1].ID, userResponses[1].ID)
}

func (suite *UserServiceTestSuite) TestUpdateUserSuccess() {
	userId := uint(1)
	userUpdateRequest := dto.UpdateUserRequest{
		Name:       "UpdatedName",
		Surname:    "UpdatedSurname",
		Patronymic: "UpdatedPatronymic",
		Address:    "UpdatedAddress",
	}

	existingUser := &models.User{
		ID:             userId,
		PassportNumber: "1234-567890",
		Surname:        "Doe",
		Name:           "John",
		Patronymic:     "Michael",
		Address:        "123 Main St",
	}

	suite.userRepoMock.EXPECT().
		GetById(userId).
		Return(existingUser, nil)

	suite.userRepoMock.EXPECT().
		Update(gomock.Any()).
		DoAndReturn(func(user *models.User) error {
			existingUser.Name = userUpdateRequest.Name
			existingUser.Surname = userUpdateRequest.Surname
			existingUser.Patronymic = userUpdateRequest.Patronymic
			existingUser.Address = userUpdateRequest.Address
			return nil
		})

	updatedUserResponse, err := suite.userService.UpdateUser(userId, userUpdateRequest)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), updatedUserResponse)
	assert.Equal(suite.T(), userId, updatedUserResponse.ID)
	assert.Equal(suite.T(), userUpdateRequest.Name, updatedUserResponse.Name)
	assert.Equal(suite.T(), userUpdateRequest.Surname, updatedUserResponse.Surname)
	assert.Equal(suite.T(), userUpdateRequest.Patronymic, updatedUserResponse.Patronymic)
	assert.Equal(suite.T(), userUpdateRequest.Address, updatedUserResponse.Address)
}

func (suite *UserServiceTestSuite) TestUpdateUserNotFound() {
	userId := uint(1)
	userUpdateRequest := dto.UpdateUserRequest{
		Name:       "UpdatedName",
		Surname:    "UpdatedSurname",
		Patronymic: "UpdatedPatronymic",
		Address:    "UpdatedAddress",
	}

	suite.userRepoMock.EXPECT().
		GetById(userId).
		Return(nil, errors.New("user not found"))

	updatedUserResponse, err := suite.userService.UpdateUser(userId, userUpdateRequest)
	assert.NotNil(suite.T(), err)
	assert.Nil(suite.T(), updatedUserResponse)
}

func (suite *UserServiceTestSuite) TestDeleteUserSuccess() {
	userId := uint(1)

	suite.userRepoMock.EXPECT().
		Delete(userId).
		Return(nil)

	err := suite.userService.DeleteUser(userId)
	assert.Nil(suite.T(), err)
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
