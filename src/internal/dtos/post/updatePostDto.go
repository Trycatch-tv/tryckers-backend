package dtos

import (
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/google/uuid"
)

type UpdatePostDto struct {
	ID      uuid.UUID        `json:"id" binding:"required"`
	Title   string           `json:"title" binding:"required"`
	Content string           `json:"content" binding:"required" `
	Image   string           `json:"image"`
	Type    enums.PostType   `json:"type"`
	Tags    string           `json:"tags"`
	Status  enums.PostStatus `json:"status" binding:"required"`
}
