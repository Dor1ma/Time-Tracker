package dto

type TaskResponse struct {
	TaskName  string `json:"task_name"`
	Hours     int    `json:"hours"`
	Minutes   int    `json:"minutes"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
