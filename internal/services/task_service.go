package services

import "github.com/Dor1ma/Time-Tracker/internal/dto"

type TaskService interface {
	StartTask(request dto.StartTaskRequest) (*dto.TaskResponse, error)
	StopTask(request dto.StopTaskRequest) (*dto.TaskResponse, error)
	GetUserTasks(userID uint, startDate string, endDate string) ([]dto.TaskResponse, error)
}
