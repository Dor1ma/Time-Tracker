package dto

type UserResponse struct {
	ID             uint   `json:"id"`
	PassportNumber string `json:"passport_number"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
}
