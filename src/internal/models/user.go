package models

import (
	"time"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name           string         `gorm:"size:100" json:"name"`
	Email          string         `gorm:"unique;size:100" json:"email"`
	Password       string         `json:"-"`
	BirthDate      *time.Time     `json:"birth_date"`
	ProfilePicture string         `json:"profile_picture"`
	GithubURL      string         `json:"github_url"`
	LinkedinURL    string         `json:"linkedin_url"`
	PitchVideo     string         `json:"pitch_video"`
	Headline       string         `json:"headline"`
	Bio            string         `json:"bio"`
	Seniority      string         `json:"seniority"`
	EnglishLevel   string         `json:"english_level"`
	EFSetScore     string         `json:"efset_score"`
	Points         int            `json:"points"`
	Role           enums.UserRole `json:"role"`
	Country        enums.Country  `json:"country"`
	Availability   string         `json:"availability"`
	Interests      string         `json:"interests"`
	Status         bool           `json:"status" default:"true"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
