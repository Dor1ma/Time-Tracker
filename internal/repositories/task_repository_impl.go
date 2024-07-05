package repositories

import (
	"github.com/Dor1ma/Time-Tracker/internal/models"
	"gorm.io/gorm"
	"time"
)

type TaskRepositoryImpl struct {
	db *gorm.DB
}

func NewTaskRepositoryImpl(db *gorm.DB) *TaskRepositoryImpl {
	return &TaskRepositoryImpl{db: db}
}

func (r *TaskRepositoryImpl) StartTask(userID uint, taskName string) (*models.Task, error) {
	task := &models.Task{
		UserID:    userID,
		TaskName:  taskName,
		StartTime: time.Now(),
	}

	if err := r.db.Create(task).Error; err != nil {
		return nil, err
	}

	return task, nil
}

func (r *TaskRepositoryImpl) StopTask(taskID uint) (*models.Task, error) {
	var task models.Task
	err := r.db.First(&task, taskID).Error
	if err != nil {
		return nil, err
	}

	task.EndTime = time.Now()
	duration := task.EndTime.Sub(task.StartTime)
	task.Hours = int(duration.Hours())
	task.Minutes = int(duration.Minutes())

	if err := r.db.Save(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *TaskRepositoryImpl) GetUserTasks(userID uint, startDate time.Time, endDate time.Time) (*[]models.Task, error) {
	var tasks []models.Task
	err := r.db.Where("user_id = ? AND start_time >= ? AND end_time <= ?", userID, startDate, endDate).
		Order("hours DESC, minutes DESC").Find(&tasks).Error
	return &tasks, err
}
