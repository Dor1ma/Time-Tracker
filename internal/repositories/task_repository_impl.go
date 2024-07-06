package repositories

import (
	"time"

	"github.com/Dor1ma/Time-Tracker/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TaskRepositoryImpl struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewTaskRepositoryImpl(db *gorm.DB, logger *logrus.Logger) *TaskRepositoryImpl {
	return &TaskRepositoryImpl{
		db:     db,
		logger: logger,
	}
}

func (r *TaskRepositoryImpl) StartTask(userID uint, taskName string) (*models.Task, error) {
	task := &models.Task{
		UserID:    userID,
		TaskName:  taskName,
		StartTime: time.Now(),
	}

	r.logger.Infof("StartTask: start adding task to database for user ID: %d, task name: %s", userID, taskName)
	if err := r.db.Create(task).Error; err != nil {
		r.logger.Debugf("StartTask: failed to add task to database: %v", err)
		return nil, err
	}

	r.logger.Infof("StartTask: successfully added task to database with ID: %d", task.ID)
	return task, nil
}

func (r *TaskRepositoryImpl) StopTask(taskID uint) (*models.Task, error) {
	var task models.Task
	if err := r.db.First(&task, taskID).Error; err != nil {
		r.logger.Debugf("StopTask: failed to find task with ID %d: %v in database", taskID, err)
		return nil, err
	}

	task.EndTime = time.Now()
	duration := task.EndTime.Sub(task.StartTime)
	task.Hours = int(duration.Hours())
	task.Minutes = int(duration.Minutes())

	r.logger.Infof("StopTask: stopping task with ID: %d", task.ID)
	if err := r.db.Save(&task).Error; err != nil {
		r.logger.Errorf("StopTask: failed to stop task: %v", err)
		return nil, err
	}

	r.logger.Infof("StopTask: successfully stopped task with ID: %d", task.ID)
	return &task, nil
}

func (r *TaskRepositoryImpl) GetUserTasks(userID uint, startDate time.Time, endDate time.Time) ([]models.Task, error) {
	var tasks []models.Task
	r.logger.Infof("GetUserTasks: fetching tasks for user ID from database: %d", userID)
	err := r.db.Where("user_id = ? AND start_time >= ? AND end_time <= ?", userID, startDate, endDate).
		Order("hours DESC, minutes DESC").Find(&tasks).Error
	if err != nil {
		r.logger.Errorf("GetUserTasks: failed to fetch tasks from database: %v", err)
		return nil, err
	}

	r.logger.Infof("GetUserTasks: successfully fetched %d tasks from database for user ID: %d", len(tasks), userID)
	return tasks, nil
}
