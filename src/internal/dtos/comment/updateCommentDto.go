package dtos

import (
	"time"

	"github.com/google/uuid"
)

type UpdateCommentDto struct {
	ID        uuid.UUID `json:"id"`
	Content   string    `json:"content"`
	Image     string    `json:"image"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    string    `json:"user_id"`
	PostId    string    `json:"post_id"`
}
