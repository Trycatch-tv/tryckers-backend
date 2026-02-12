package dtos

import (
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/google/uuid"
)

// UpdatePostDto representa los datos para actualizar un post existente
// Todos los campos son opcionales, solo se actualizarán los campos proporcionados
type UpdatePostDto struct {
	ID       uuid.UUID         `json:"id" swaggerignore:"true"`
	Title    *string           `json:"title,omitempty" binding:"omitempty,min=3,max=500" example:"Título actualizado"`
	Content  *string           `json:"content,omitempty" binding:"omitempty,min=10,max=10000" example:"Contenido actualizado..."`
	Image    *string           `json:"image,omitempty" binding:"omitempty,url" example:"https://example.com/new-image.jpg"`
	Type     *enums.PostType   `json:"type,omitempty" binding:"omitempty,oneof=regular story video" example:"regular"`
	Tags     []string          `json:"tags,omitempty" binding:"omitempty,max=10,dive,min=2,max=50"`
	Status   *enums.PostStatus `json:"status,omitempty" binding:"omitempty,oneof=published edited draft deleted" example:"published"`
	UserID   uuid.UUID         `json:"user_id" swaggerignore:"true"`
	MediaURL *string           `json:"media_url,omitempty" binding:"omitempty,url" example:"https://example.com/video.mp4"`
}

// HasChanges verifica si el DTO contiene algún cambio
func (dto *UpdatePostDto) HasChanges() bool {
	return dto.Title != nil ||
		dto.Content != nil ||
		dto.Image != nil ||
		dto.Type != nil ||
		dto.Tags != nil ||
		dto.Status != nil ||
		dto.MediaURL != nil
}

// Validate realiza validaciones adicionales personalizadas
func (dto *UpdatePostDto) Validate() error {
	if dto.Type != nil && *dto.Type == enums.VideoPost && dto.MediaURL == nil {
		return ErrVideoPostRequiresMedia
	}
	return nil
}
