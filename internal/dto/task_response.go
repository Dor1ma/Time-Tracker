package dto

type TaskResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	TaskName  string `json:"task_name"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
