package dto

type TaskResponse struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	TaskName  string `json:"task_name"`
	Hours     int    `json:"hours"`
	Minutes   int    `json:"minutes"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
