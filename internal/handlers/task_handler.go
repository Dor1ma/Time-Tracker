package handlers

import (
	"github.com/Dor1ma/Time-Tracker/internal/dto"
	"github.com/Dor1ma/Time-Tracker/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	taskService services.TaskService
	logger      *logrus.Logger
}

func NewTaskHandler(taskService services.TaskService, logger *logrus.Logger) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
		logger:      logger,
	}
}

// StartTask godoc
// @Summary Start a new task
// @Description Start a new task for a user
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body dto.StartTaskRequest true "Start task request"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /tasks/start [post]
func (h *TaskHandler) StartTask(c *gin.Context) {
	var request dto.StartTaskRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Debugf("StartTask: invalid request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Infof("StartTask: received request to start task for user ID: %d, task name: %s", request.UserID, request.TaskName)
	task, err := h.taskService.StartTask(request)
	if err != nil {
		h.logger.Debugf("StartTask: failed to start task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Infof("StartTask: successfully started task with ID: %d", task.ID)
	c.JSON(http.StatusOK, task)
}

// StopTask godoc
// @Summary Stop an existing task
// @Description Stop an existing task for a user
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body dto.StopTaskRequest true "Stop task request"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /tasks/stop [post]
func (h *TaskHandler) StopTask(c *gin.Context) {
	var request dto.StopTaskRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Debugf("StopTask: invalid task ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Infof("StopTask: received request to stop task with ID: %d", request.TaskID)
	task, err := h.taskService.StopTask(request)
	if err != nil {
		h.logger.Debugf("StopTask: failed to stop task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Infof("StopTask: successfully stopped task with ID: %d", task.ID)
	c.JSON(http.StatusOK, task)
}

// GetUserTasks godoc
// @Summary Get user tasks
// @Description Get tasks for a user within a specified date range, sorted by total time spent
// @Tags tasks
// @Accept json
// @Produce json
// @Param user_id path int true "User ID"
// @Param start_date query string true "Start date in format YYYY-MM-DD"
// @Param end_date query string true "End date in format YYYY-MM-DD"
// @Success 200 {array} models.Task
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /users/{user_id}/tasks [get]
func (h *TaskHandler) GetUserTasks(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		h.logger.Debugf("GetUserTasks: invalid user ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	h.logger.Infof("GetUserTasks: received request to fetch tasks for user ID: %d", userID)
	tasks, err := h.taskService.GetUserTasks(uint(userID), startDate, endDate)
	if err != nil {
		h.logger.Debugf("GetUserTasks: failed to fetch tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Infof("GetUserTasks: successfully fetched %d tasks for user ID: %d", len(tasks), userID)
	c.JSON(http.StatusOK, tasks)
}
