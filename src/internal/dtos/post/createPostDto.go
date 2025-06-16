package dtos

import (
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/google/uuid"
)

type CreatePostDto struct {
	// ID      uuid.UUID    `json:"id"`
	Title   string       `json:"title" binding:"required"`
	Content string       `json:"content" binding:"required"`
	Image   string       `json:"image"`
	Type    string       `json:"type"`
	Tags    string       `json:"tags"`
	Status  enums.Status `json:"status" binding:"required"`
	UserId  uuid.UUID    `json:"user_id" binding:"required"`
}
