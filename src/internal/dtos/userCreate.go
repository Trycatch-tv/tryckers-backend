package dtos

import "github.com/Trycatch-tv/tryckers-backend/src/internal/enums"

// CreateUserDTO represents the user registration request payload
type CreateUserDTO struct {
	Name     string        `json:"name" binding:"required" example:"John Doe"`
	Country  enums.Country `json:"country" binding:"required" example:"US"`
	Email    string        `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password string        `json:"password" binding:"required,min=8" example:"securepassword123"`
} // @name CreateUserDTO
