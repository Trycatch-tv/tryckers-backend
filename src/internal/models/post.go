package models

import (
	"strings"
	"time"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/google/uuid"
)

// Post representa un post en la plataforma
type Post struct {
	ID      uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Title   string         `gorm:"size:500;not null" json:"title"`
	Content string         `gorm:"type:text;not null" json:"content"`
	Image   string         `gorm:"size:1000" json:"image"`
	Type    enums.PostType `gorm:"size:50;default:'regular'" json:"type"`
	// Tags es una cadena de tags separados por coma asociados al post
	Tags       string           `gorm:"size:1000" json:"tags"`
	Status     enums.PostStatus `gorm:"size:50;default:'draft'" json:"status"`
	MediaURL   string           `gorm:"size:1000" json:"media_url"`
	CreatedAt  *time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  *time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	UserID     uuid.UUID        `gorm:"type:uuid;not null;index" json:"user_id"`
	User       User             `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
	VotesCount int64            `gorm:"column:votes_count;->" json:"votes_count"`
	Comments   []Comment        `gorm:"foreignKey:PostID;references:ID" json:"comments,omitempty"`
}

// TableName especifica el nombre de la tabla
func (Post) TableName() string {
	return "posts"
}

// BeforeCreate hook que se ejecuta antes de crear un post
func (p *Post) BeforeCreate(tx interface{}) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	if p.Type == "" {
		p.Type = enums.RegularPost
	}
	if p.Status == "" {
		p.Status = enums.DRAFT
	}
	return nil
}

// GetTagsSlice retorna los tags como un slice de strings
func (p *Post) GetTagsSlice() []string {
	if p.Tags == "" {
		return []string{}
	}
	tags := strings.Split(p.Tags, ",")
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// SetTagsFromSlice establece los tags desde un slice de strings
func (p *Post) SetTagsFromSlice(tags []string) {
	cleanTags := make([]string, 0, len(tags))
	for _, tag := range tags {
		trimmed := strings.TrimSpace(tag)
		if trimmed != "" {
			cleanTags = append(cleanTags, trimmed)
		}
	}
	p.Tags = strings.Join(cleanTags, ",")
}

// IsPublished verifica si el post está publicado
func (p *Post) IsPublished() bool {
	return p.Status == enums.PUBLISHED
}

// IsDraft verifica si el post es un borrador
func (p *Post) IsDraft() bool {
	return p.Status == enums.DRAFT
}

// IsDeleted verifica si el post está eliminado
func (p *Post) IsDeleted() bool {
	return p.Status == enums.DELETED
}

// IsVideoPost verifica si el post es de tipo video
func (p *Post) IsVideoPost() bool {
	return p.Type == enums.VideoPost
}

// CanBeEditedBy verifica si un usuario puede editar el post
func (p *Post) CanBeEditedBy(userID uuid.UUID) bool {
	return p.UserID == userID
}

// CanBeDeletedBy verifica si un usuario puede eliminar el post
func (p *Post) CanBeDeletedBy(userID uuid.UUID, isAdmin bool) bool {
	return p.UserID == userID || isAdmin
}
