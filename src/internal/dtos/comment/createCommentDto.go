package dtos

import (
	"github.com/google/uuid"
)

type CreateCommentDto struct {
	// ID      uuid.UUID    `json:"id"`
	Content string `json:"content" binding:"required"`
	// Status  enums.Status `json:"status" binding:"required"`
	Image  string    `json:"image"`
	Status bool      `json:"status"`
	UserId uuid.UUID `json:"user_id" binding:"required"`
	PostId uuid.UUID `json:"post_id" binding:"required"`
}
