package dtos

import (
	"github.com/google/uuid"
)

type UpdateCommentDto struct {
	ID      uuid.UUID `json:"id"`
	Content string    `json:"content"`
}
