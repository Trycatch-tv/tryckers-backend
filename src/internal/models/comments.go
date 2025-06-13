package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	IdPublic  string    `gorm:"size:1000" json:"content"`
	Content   string    `gorm:"size:1000" json:"content"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserId    uuid.UUID `json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;references:ID" json:"user"`
	PostId    uuid.UUID `json:"post_id"`
	Post      Post      `gorm:"foreignKey:PostID;references:ID" json:"post"`
}
