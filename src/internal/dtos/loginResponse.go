package dtos

import "github.com/Trycatch-tv/tryckers-backend/src/internal/models"

type LoginResponse struct {
	UserData     models.User `json:"user_data"`
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
}

// RefreshTokenRequest representa la solicitud para refrescar el token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse representa la respuesta con los nuevos tokens
type RefreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
