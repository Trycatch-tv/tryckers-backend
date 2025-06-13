package dtos

import (
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/google/uuid"
)

type CreatePostRequest struct {
	ID      uuid.UUID    `json:"id"`
	Title   string       `json:"title" binding:"required"`
	Content string       `json:"content" binding:"required"`
	Status  enums.Status `json:"status" binding:"required"`
	UserId  uuid.UUID    `json:"user_id" binding:"required"`
}
