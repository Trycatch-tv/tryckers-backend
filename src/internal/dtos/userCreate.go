package dtos

import "github.com/Trycatch-tv/tryckers-backend/src/internal/enums"

type CreateUserDTO struct {
	Name     string        `json:"name" binding:"required"`
	Country  enums.Country `json:"country" binding:"required" `
	Email    string        `json:"email" binding:"required,email"`
	Password string        `json:"password" binding:"required,min=8"`
}
