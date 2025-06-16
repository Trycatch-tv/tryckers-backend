package dtos

import (
	"time"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
)

type UpdatePostDto struct {
	ID        string       `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Image     string       `json:"image"`
	Type      enums.Type   `json:"type"`
	Tags      string       `json:"tags"`
	Status    enums.Status `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	UserId    string       `json:"user_id"`
}
