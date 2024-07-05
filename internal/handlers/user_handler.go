package handlers

import (
	"github.com/Dor1ma/Time-Tracker/internal/dto"
	"github.com/Dor1ma/Time-Tracker/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userService services.UserService
	logger      *logrus.Logger
}

func NewUserHandler(userService services.UserService, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Debugf("CreateUser: invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Infof("CreateUser: creating user with passport number: %s", req.PassportNumber)
	user, err := h.userService.CreateUser(req.PassportNumber)
	if err != nil {
		h.logger.Debugf("CreateUser: failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Infof("CreateUser: user created with ID: %d", user.ID)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Debugf("UpdateUser: invalid user ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var userUpdateRequest dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&userUpdateRequest); err != nil {
		h.logger.Debugf("UpdateUser: invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Infof("UpdateUser: updating user with ID: %d", userID)
	user, err := h.userService.UpdateUser(uint(userID), userUpdateRequest)
	if err != nil {
		h.logger.Debugf("UpdateUser: failed to update user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Infof("UpdateUser: user updated with ID: %d", userID)
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Debugf("DeleteUser: invalid user ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	h.logger.Infof("DeleteUser: deleting user with ID: %d", userID)
	err = h.userService.DeleteUser(uint(userID))
	if err != nil {
		h.logger.Debugf("DeleteUser: failed to delete user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Infof("DeleteUser: user deleted with ID: %d", userID)
	c.JSON(http.StatusOK, gin.H{"user": nil})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	h.logger.Info("GetUsers: fetching all users")
	filters := make(map[string]interface{})
	for key, values := range c.Request.URL.Query() {
		filters[key] = values[0]
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	users, err := h.userService.GetUsersWithFiltersAndPagination(filters, page, pageSize)
	if err != nil {
		h.logger.Debugf("GetUsers: failed to fetch users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Infof("GetUsers: fetched %d users", len(users))
	c.JSON(http.StatusOK, users)
}
