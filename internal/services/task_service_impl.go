package services

import (
	"github.com/Dor1ma/Time-Tracker/internal/models"
	"github.com/Dor1ma/Time-Tracker/internal/repositories"
	"time"
)

type TaskServiceImpl struct {
	taskRepo repositories.TaskRepository
}

func NewTaskServiceImpl(taskRepo repositories.TaskRepository) *TaskServiceImpl {
	return &TaskServiceImpl{taskRepo: taskRepo}
}

func (s *TaskServiceImpl) StartTask(userID uint, taskName string) (*models.Task, error) {
	return s.taskRepo.StartTask(userID, taskName)
}

func (s *TaskServiceImpl) StopTask(taskID uint) (*models.Task, error) {
	return s.taskRepo.StopTask(taskID)
}

func (s *TaskServiceImpl) GetUserTasks(userID uint, startDate string, endDate string) ([]models.Task, error) {
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return nil, err
	}

	return s.taskRepo.GetUserTasks(userID, start, end)
}
