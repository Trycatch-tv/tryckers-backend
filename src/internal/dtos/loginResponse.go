package dtos

import "github.com/Trycatch-tv/tryckers-backend/src/internal/models"

type LoginResponse struct {
	UserData models.User
	Token    string
}
