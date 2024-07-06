package dto

type UpdateUserRequest struct {
	Surname    string `json:"surname" binding:"required"`
	Name       string `json:"name" binding:"required"`
	Patronymic string `json:"patronymic" binding:"required"`
	Address    string `json:"address" binding:"required"`
}
