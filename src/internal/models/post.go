package models

import (
	"time"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/google/uuid"
)

type Post struct {
	ID      uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Title   string         `gorm:"size:1000;not null" json:"title"`
	Content string         `gorm:"size:1000;not null" json:"content"`
	Image   string         `gorm:"size:1000" json:"image"`
	Type    enums.PostType `gorm:"size:50;default:'text'" json:"type"`
	// Tags is a comma-separated string of tags associated with the post
	Tags      string           `gorm:"size:1000" json:"tags"`
	Votes     int              `gorm:"default:0" json:"votes"`
	Status    enums.PostStatus `gorm:"size:50;default:'draft'" json:"status"`
	CreatedAt *time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	UserID    uuid.UUID        `gorm:"type:uuid;not null" json:"user_id"`
	User      User             `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
}
