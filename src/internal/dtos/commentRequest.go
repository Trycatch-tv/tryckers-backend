package dtos

import (
	"github.com/google/uuid"
)

type CreateCommentRequest struct {
	// ID      uuid.UUID    `json:"id"`
	Content string `json:"content" binding:"required"`
	// Status  enums.Status `json:"status" binding:"required"`
	UserId uuid.UUID `json:"user_id" binding:"required"`
	PostId uuid.UUID `json:"post_id" binding:"required"`
}
