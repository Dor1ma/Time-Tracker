package dto

type StartTaskRequest struct {
	UserID   uint   `json:"user_id" binding:"required"`
	TaskName string `json:"task_name" binding:"required"`
}
