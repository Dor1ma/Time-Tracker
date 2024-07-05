package services

import (
	"github.com/Dor1ma/Time-Tracker/internal/dto"
	"github.com/Dor1ma/Time-Tracker/internal/repositories"
	"time"
)

type TaskServiceImpl struct {
	taskRepo repositories.TaskRepository
}

func NewTaskServiceImpl(taskRepo repositories.TaskRepository) *TaskServiceImpl {
	return &TaskServiceImpl{taskRepo: taskRepo}
}

func (s *TaskServiceImpl) StartTask(request dto.StartTaskRequest) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.StartTask(request.UserID, request.TaskName)
	if err != nil {
		return nil, err
	}

	return &dto.TaskResponse{
		ID:        task.ID,
		UserID:    task.UserID,
		TaskName:  task.TaskName,
		StartTime: task.StartTime.Format(time.RFC3339),
	}, nil
}

func (s *TaskServiceImpl) StopTask(request dto.StopTaskRequest) (*dto.TaskResponse, error) {
	task, err := s.taskRepo.StopTask(request.TaskID)
	if err != nil {
		return nil, err
	}

	return &dto.TaskResponse{
		ID:        task.ID,
		UserID:    task.UserID,
		TaskName:  task.TaskName,
		StartTime: task.StartTime.Format(time.RFC3339),
		EndTime:   task.EndTime.Format(time.RFC3339),
	}, nil
}

func (s *TaskServiceImpl) GetUserTasks(userID uint, startDate string, endDate string) ([]dto.TaskResponse, error) {
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return nil, err
	}

	tasks, err := s.taskRepo.GetUserTasks(userID, start, end)
	if err != nil {
		return nil, err
	}

	taskResponses := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = dto.TaskResponse{
			ID:        task.ID,
			UserID:    task.UserID,
			TaskName:  task.TaskName,
			StartTime: task.StartTime.Format(time.RFC3339),
			EndTime:   task.EndTime.Format(time.RFC3339),
		}
	}

	return taskResponses, nil
}
