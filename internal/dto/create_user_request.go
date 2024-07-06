package dto

type CreateUserRequest struct {
	PassportNumber string `json:"passportNumber" binding:"required"`
}
