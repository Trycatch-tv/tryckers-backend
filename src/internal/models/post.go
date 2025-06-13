package models

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	IdPublic  string     `gorm:"size:1000" json:"id_public"`
	Title     string     `gorm:"size:1000;not null" json:"title"`
	Content   string     `gorm:"size:1000;not null" json:"content"`
	Status    string     `gorm:"size:50;default:'draft'" json:"status"`
	CreatedAt *time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	User      User       `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`
}
