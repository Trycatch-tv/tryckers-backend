package dtos

import (
	"time"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
)

// Type represents the type of post
var _ = enums.PostType(0) // Ensure the Type enum is imported

type ResponsePostDto struct {
	ID        string           `json:"id"`
	Title     string           `json:"title"`
	Content   string           `json:"content"`
	Image     string           `json:"image"`
	Type      enums.PostType   `json:"string"`
	Tags      string           `json:"tags"`
	Status    enums.PostStatus `json:"status"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
	UserId    string           `json:"user_id"`
}
