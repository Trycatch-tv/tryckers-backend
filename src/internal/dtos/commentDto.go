package dtos

import (
	"time"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
)

type CommentDto struct {
	ID        string       `json:"id"`
	Content   string       `json:"content"`
	Status    enums.Status `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	UserId    string       `json:"user_id"`
	PostId    string       `json:"post_id"`
}
