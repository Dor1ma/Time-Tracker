package tests

import (
	"errors"
	"github.com/Dor1ma/Time-Tracker/internal/repositories"
	"github.com/golang/mock/gomock"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Dor1ma/Time-Tracker/internal/dto"
	"github.com/Dor1ma/Time-Tracker/internal/models"
	"github.com/Dor1ma/Time-Tracker/internal/services"
	"github.com/sirupsen/logrus"
)

func TestStartTask_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositories.NewMockTaskRepository(ctrl)
	logger := logrus.New()

	service := services.NewTaskServiceImpl(mockRepo, logger)

	userID := uint(1)
	taskName := "Sample Task"
	request := dto.StartTaskRequest{UserID: userID, TaskName: taskName}

	expectedTask := &models.Task{
		ID:        1,
		UserID:    userID,
		TaskName:  taskName,
		StartTime: time.Now(),
	}

	mockRepo.EXPECT().StartTask(userID, taskName).Return(expectedTask, nil)

	taskResponse, err := service.StartTask(request)

	assert.NoError(t, err)
	assert.NotNil(t, taskResponse)
	assert.Equal(t, expectedTask.ID, taskResponse.ID)
	assert.Equal(t, expectedTask.UserID, taskResponse.UserID)
	assert.Equal(t, expectedTask.TaskName, taskResponse.TaskName)
	assert.Equal(t, expectedTask.StartTime.Format(time.RFC3339), taskResponse.StartTime)
}

func TestStartTask_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositories.NewMockTaskRepository(ctrl)
	logger := logrus.New()

	service := services.NewTaskServiceImpl(mockRepo, logger)

	userID := uint(1)
	taskName := "Sample Task"
	request := dto.StartTaskRequest{UserID: userID, TaskName: taskName}

	expectedError := errors.New("repository error")

	mockRepo.EXPECT().StartTask(userID, taskName).Return(nil, expectedError)

	taskResponse, err := service.StartTask(request)

	assert.Error(t, err)
	assert.Nil(t, taskResponse)
	assert.Equal(t, expectedError, err)
}

func TestStopTask_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositories.NewMockTaskRepository(ctrl)
	logger := logrus.New()

	service := services.NewTaskServiceImpl(mockRepo, logger)

	taskID := uint(1)

	expectedTask := &models.Task{
		ID:        taskID,
		UserID:    1,
		TaskName:  "Sample Task",
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
	}

	mockRepo.EXPECT().StopTask(taskID).Return(expectedTask, nil)

	request := dto.StopTaskRequest{TaskID: taskID}
	taskResponse, err := service.StopTask(request)

	assert.NoError(t, err)
	assert.NotNil(t, taskResponse)
	assert.Equal(t, expectedTask.ID, taskResponse.ID)
	assert.Equal(t, expectedTask.UserID, taskResponse.UserID)
	assert.Equal(t, expectedTask.TaskName, taskResponse.TaskName)
	assert.Equal(t, expectedTask.StartTime.Format(time.RFC3339), taskResponse.StartTime)
	assert.Equal(t, expectedTask.EndTime.Format(time.RFC3339), taskResponse.EndTime)
}

func TestStopTask_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositories.NewMockTaskRepository(ctrl)
	logger := logrus.New()

	service := services.NewTaskServiceImpl(mockRepo, logger)

	taskID := uint(1)

	expectedError := errors.New("repository error")

	mockRepo.EXPECT().StopTask(taskID).Return(nil, expectedError)

	request := dto.StopTaskRequest{TaskID: taskID}
	taskResponse, err := service.StopTask(request)

	assert.Error(t, err)
	assert.Nil(t, taskResponse)
	assert.Equal(t, expectedError, err)
}

func TestGetUserTasks_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositories.NewMockTaskRepository(ctrl)
	logger := logrus.New()

	service := services.NewTaskServiceImpl(mockRepo, logger)

	userID := uint(1)
	startDate := "2024-07-01"
	endDate := "2024-07-31"

	expectedTasks := []models.Task{
		{
			ID:        1,
			UserID:    userID,
			TaskName:  "Task 1",
			Hours:     2,
			Minutes:   30,
			StartTime: time.Now(),
			EndTime:   time.Now().Add(time.Hour),
		},
		{
			ID:        2,
			UserID:    userID,
			TaskName:  "Task 2",
			Hours:     1,
			Minutes:   45,
			StartTime: time.Now().Add(-24 * time.Hour),
			EndTime:   time.Now().Add(-23 * time.Hour).Add(45 * time.Minute),
		},
	}

	startTime, _ := time.Parse("2006-01-02", startDate)
	endTime, _ := time.Parse("2006-01-02", endDate)

	mockRepo.EXPECT().GetUserTasks(userID, startTime, endTime).Return(expectedTasks, nil)

	taskResponses, err := service.GetUserTasks(userID, startDate, endDate)

	assert.NoError(t, err)
	assert.NotNil(t, taskResponses)
	assert.Equal(t, len(expectedTasks), len(taskResponses))

	for i, expectedTask := range expectedTasks {
		assert.Equal(t, expectedTask.ID, taskResponses[i].ID)
		assert.Equal(t, expectedTask.UserID, taskResponses[i].UserID)
		assert.Equal(t, expectedTask.TaskName, taskResponses[i].TaskName)
		assert.Equal(t, expectedTask.Hours, taskResponses[i].Hours)
		assert.Equal(t, expectedTask.Minutes, taskResponses[i].Minutes)
		assert.Equal(t, expectedTask.StartTime.Format(time.RFC3339), taskResponses[i].StartTime)
		assert.Equal(t, expectedTask.EndTime.Format(time.RFC3339), taskResponses[i].EndTime)
	}
}

func TestGetUserTasks_InvalidStartDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositories.NewMockTaskRepository(ctrl)
	logger := logrus.New()

	service := services.NewTaskServiceImpl(mockRepo, logger)

	userID := uint(1)
	startDate := "invalid-date"
	endDate := "2024-07-31"

	taskResponses, err := service.GetUserTasks(userID, startDate, endDate)

	assert.Error(t, err)
	assert.Nil(t, taskResponses)
	assert.Contains(t, err.Error(), "parsing time")
}

func TestGetUserTasks_InvalidEndDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repositories.NewMockTaskRepository(ctrl)
	logger := logrus.New()

	service := services.NewTaskServiceImpl(mockRepo, logger)

	userID := uint(1)
	startDate := "2024-07-01"
	endDate := "invalid-date"

	taskResponses, err := service.GetUserTasks(userID, startDate, endDate)

	assert.Error(t, err)
	assert.Nil(t, taskResponses)
	assert.Contains(t, err.Error(), "parsing time")
}
