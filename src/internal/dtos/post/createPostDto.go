package dtos

import (
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/google/uuid"
)

// CreatePostDto representa los datos necesarios para crear un nuevo post
type CreatePostDto struct {
	Title    string           `json:"title" binding:"required,min=3,max=500" example:"Mi primer post"`
	Content  string           `json:"content" binding:"required,min=10,max=10000" example:"Este es el contenido de mi post..."`
	Image    string           `json:"image" binding:"omitempty,url" example:"https://example.com/image.jpg"`
	Type     enums.PostType   `json:"type" binding:"required,oneof=regular story video" example:"regular"`
	Tags     []string         `json:"tags" binding:"omitempty,max=10,dive,min=2,max=50" example:"go,backend,api"`
	Status   enums.PostStatus `json:"status" binding:"required,oneof=published draft" example:"draft"`
	UserID   uuid.UUID        `json:"user_id" binding:"required" swaggerignore:"true"`
	MediaURL string           `json:"media_url" binding:"omitempty,url" example:"https://example.com/video.mp4"`
}

// Validate realiza validaciones adicionales personalizadas
func (dto *CreatePostDto) Validate() error {
	// Validar que posts de tipo video tengan media_url
	if dto.Type == enums.VideoPost && dto.MediaURL == "" {
		return ErrVideoPostRequiresMedia
	}
	return nil
}
