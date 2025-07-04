package dtos

import (
	"github.com/google/uuid"
)

type CreateCommentDto struct {
	Content string    `json:"content" binding:"required"`
	Image   string    `json:"image"`
	Status  bool      `json:"status"`
	UserId  uuid.UUID `json:"user_id" binding:"required"`
	PostId  uuid.UUID `json:"post_id" binding:"required"`
}
