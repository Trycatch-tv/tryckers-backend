package dtos

import (
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/google/uuid"
)

type UpdatePostDto struct {
	ID      uuid.UUID        `json:"id"`
	Title   string           `json:"title"`
	Content string           `json:"content"`
	Image   string           `json:"image"`
	Type    enums.PostType   `json:"type"`
	Tags    string           `json:"tags"`
	Status  enums.PostStatus `json:"status"`
	UserId  string           `json:"user_id"`
}
