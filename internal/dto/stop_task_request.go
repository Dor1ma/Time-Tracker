package dto

type StopTaskRequest struct {
	TaskID uint `json:"task_id" binding:"required"`
}
