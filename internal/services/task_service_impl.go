package services

import (
	"github.com/Dor1ma/Time-Tracker/internal/dto"
	"github.com/Dor1ma/Time-Tracker/internal/repositories"
	"github.com/sirupsen/logrus"
	"time"
)

type TaskServiceImpl struct {
	taskRepo repositories.TaskRepository
	logger   *logrus.Logger
}

func NewTaskServiceImpl(taskRepo repositories.TaskRepository, logger *logrus.Logger) *TaskServiceImpl {
	return &TaskServiceImpl{
		taskRepo: taskRepo,
		logger:   logger,
	}
}

func (s *TaskServiceImpl) StartTask(request dto.StartTaskRequest) (*dto.TaskResponse, error) {
	s.logger.Infof("StartTask: starting task for user ID: %d, task name: %s", request.UserID, request.TaskName)
	task, err := s.taskRepo.StartTask(request.UserID, request.TaskName)
	if err != nil {
		s.logger.Debugf("StartTask: failed to start task: %v", err)
		return nil, err
	}

	s.logger.Infof("StartTask: task started with ID: %d", task.ID)
	return &dto.TaskResponse{
		ID:        task.ID,
		UserID:    task.UserID,
		TaskName:  task.TaskName,
		StartTime: task.StartTime.Format(time.RFC3339),
	}, nil
}

func (s *TaskServiceImpl) StopTask(request dto.StopTaskRequest) (*dto.TaskResponse, error) {
	s.logger.Infof("StopTask: stopping task with ID: %d", request.TaskID)
	task, err := s.taskRepo.StopTask(request.TaskID)
	if err != nil {
		s.logger.Debug("StopTask: failed to stop task: %v", err)
		return nil, err
	}

	s.logger.Infof("StopTask: task stopped with ID: %d", task.ID)
	return &dto.TaskResponse{
		ID:        task.ID,
		UserID:    task.UserID,
		TaskName:  task.TaskName,
		StartTime: task.StartTime.Format(time.RFC3339),
		EndTime:   task.EndTime.Format(time.RFC3339),
	}, nil
}

func (s *TaskServiceImpl) GetUserTasks(userID uint, startDate string, endDate string) ([]dto.TaskResponse, error) {
	s.logger.Infof("GetUserTasks: fetching tasks for user ID: %d, start date: %s, end date: %s", userID, startDate, endDate)
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		s.logger.Debugf("GetUserTasks: invalid start date: %v", err)
		return nil, err
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		s.logger.Debugf("GetUserTasks: invalid end date: %v", err)
		return nil, err
	}

	tasks, err := s.taskRepo.GetUserTasks(userID, start, end)
	if err != nil {
		s.logger.Debugf("GetUserTasks: failed to fetch tasks: %v", err)
		return nil, err
	}

	taskResponses := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = dto.TaskResponse{
			ID:        task.ID,
			UserID:    task.UserID,
			TaskName:  task.TaskName,
			Hours:     task.Hours,
			Minutes:   task.Minutes,
			StartTime: task.StartTime.Format(time.RFC3339),
			EndTime:   task.EndTime.Format(time.RFC3339),
		}
	}

	s.logger.Infof("GetUserTasks: fetched %d tasks for user ID: %d", len(tasks), userID)
	return taskResponses, nil
}
