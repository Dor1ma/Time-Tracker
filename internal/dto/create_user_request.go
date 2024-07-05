package dto

type CreateUserRequest struct {
	PassportNumber string `json:"passport_number" binding:"required"`
}
