package dtos

import (
	"strings"
	"time"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
)

// ResponsePostDto representa la respuesta de un post para el cliente
type ResponsePostDto struct {
	ID             string           `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title          string           `json:"title" example:"Mi primer post"`
	Content        string           `json:"content" example:"Este es el contenido de mi post..."`
	ContentPreview string           `json:"content_preview,omitempty" example:"Este es el contenido..."`
	Image          string           `json:"image,omitempty" example:"https://example.com/image.jpg"`
	Type           enums.PostType   `json:"type" example:"regular"`
	Tags           []string         `json:"tags" example:"go,backend,api"`
	Status         enums.PostStatus `json:"status" example:"published"`
	MediaURL       string           `json:"media_url,omitempty" example:"https://example.com/video.mp4"`
	CreatedAt      time.Time        `json:"created_at" example:"2026-02-11T10:00:00Z"`
	UpdatedAt      time.Time        `json:"updated_at" example:"2026-02-11T12:30:00Z"`
	UserID         string           `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Author         *AuthorDto       `json:"author,omitempty"`
	UserVote       int8             `json:"user_vote" example:"1"`
	VotesCount     int64            `json:"votes_count" example:"42"`
	CommentsCount  int64            `json:"comments_count" example:"5"`
	IsEdited       bool             `json:"is_edited" example:"false"`
}

// AuthorDto representa información básica del autor de un post
type AuthorDto struct {
	ID        string `json:"id" example:"550e8400-e29b-41d4-a716-446655440001"`
	Username  string `json:"username" example:"johndoe"`
	AvatarURL string `json:"avatar_url,omitempty" example:"https://example.com/avatar.jpg"`
}

// ResponsePostListDto representa una lista paginada de posts
type ResponsePostListDto struct {
	Posts      []ResponsePostDto `json:"posts"`
	Total      int64             `json:"total" example:"100"`
	Page       int               `json:"page" example:"1"`
	PageSize   int               `json:"page_size" example:"10"`
	TotalPages int               `json:"total_pages" example:"10"`
	HasNext    bool              `json:"has_next" example:"true"`
	HasPrev    bool              `json:"has_prev" example:"false"`
}

// ToResponsePostDto convierte un modelo Post a ResponsePostDto
func ToResponsePostDto(post *models.Post, userVote int8) ResponsePostDto {
	dto := ResponsePostDto{
		ID:         post.ID.String(),
		Title:      post.Title,
		Content:    post.Content,
		Image:      post.Image,
		Type:       post.Type,
		Status:     post.Status,
		MediaURL:   post.MediaURL,
		VotesCount: post.VotesCount,
		UserVote:   userVote,
		UserID:     post.UserID.String(),
		IsEdited:   post.Status == enums.EDITED,
	}

	// Parsear tags
	if post.Tags != "" {
		dto.Tags = parseTagsFromString(post.Tags)
	}

	// Agregar timestamps
	if post.CreatedAt != nil {
		dto.CreatedAt = *post.CreatedAt
	}
	if post.UpdatedAt != nil {
		dto.UpdatedAt = *post.UpdatedAt
	}

	// Generar preview del contenido (primeros 150 caracteres)
	if len(post.Content) > 150 {
		dto.ContentPreview = post.Content[:150] + "..."
	}

	// Agregar información del autor si está disponible
	if post.User.ID.String() != "00000000-0000-0000-0000-000000000000" {
		dto.Author = &AuthorDto{
			ID:        post.User.ID.String(),
			Username:  post.User.Username,
			AvatarURL: post.User.ProfilePicture,
		}
	}

	return dto
}

// ToResponsePostListDto convierte una lista de posts a ResponsePostListDto
func ToResponsePostListDto(posts []models.Post, total int64, page, pageSize int, userVotes map[string]int8) ResponsePostListDto {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	dtos := make([]ResponsePostDto, len(posts))
	for i, post := range posts {
		vote := int8(0)
		if userVotes != nil {
			if v, ok := userVotes[post.ID.String()]; ok {
				vote = v
			}
		}
		dtos[i] = ToResponsePostDto(&post, vote)
	}

	return ResponsePostListDto{
		Posts:      dtos,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
	}
}

// parseTagsFromString convierte una cadena de tags separados por coma a un slice
func parseTagsFromString(tags string) []string {
	if tags == "" {
		return []string{}
	}
	result := []string{}
	for _, tag := range splitAndTrim(tags, ",") {
		if tag != "" {
			result = append(result, tag)
		}
	}
	return result
}

// splitAndTrim divide una cadena y elimina espacios en blanco
func splitAndTrim(s, sep string) []string {
	parts := make([]string, 0)
	for _, part := range strings.Split(s, sep) {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}
	return parts
}
