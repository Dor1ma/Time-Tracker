package repositories

import (
	"github.com/Dor1ma/Time-Tracker/internal/models"
	"time"
)

type TaskRepository interface {
	StartTask(userID uint, taskName string) (*models.Task, error)
	StopTask(taskID uint) (*models.Task, error)
	GetUserTasks(userID uint, startDate, endDate time.Time) (*[]models.Task, error)
}
