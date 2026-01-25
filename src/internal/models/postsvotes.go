package models

import (
	"time"

	"github.com/google/uuid"
)

type PostVote struct {
	PostID    uuid.UUID `gorm:"type:uuid;primaryKey;autoIncrement:false"`
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey;autoIncrement:false"`
	Vote      int8      `gorm:"not null;default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
	User      User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"`
	Post      Post `gorm:"foreignKey:PostID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"post"`
}
