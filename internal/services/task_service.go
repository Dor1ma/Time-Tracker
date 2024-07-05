package services

import "github.com/Dor1ma/Time-Tracker/internal/models"

type TaskService interface {
	StartTask(userID uint, taskName string) (*models.Task, error)
	StopTask(taskID uint) (*models.Task, error)
	GetUserTasks(userID uint, startDate string, endDate string) ([]models.Task, error)
}
