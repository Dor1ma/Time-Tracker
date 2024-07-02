package models

import "time"

type User struct {
	ID             uint   `gorm:"primary_key"`
	PassportNumber string `gorm:"unique; not null"`
	Surname        string `gorm:"not null"`
	Name           string `gorm:"not null"`
	Patronymic     string
	Address        string `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
