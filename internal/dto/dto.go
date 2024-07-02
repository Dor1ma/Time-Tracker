package dto

type CreateUserRequest struct {
	PassportNumber string `json:"passport_number" binding:"required"`
}

type ExternalAPIResponse struct {
	Surname    string `json:"surname"`
	Name       string `json:"name"`
	Patronymic string `json:"patronymic"`
	Address    string `json:"address"`
}
