package dtos

import (
	"time"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
)

// UpdateProfileDTO represents the payload for updating the authenticated user's profile
type UpdateProfileDTO struct {
	Name         *string        `json:"name,omitempty" example:"John Doe"`
	Headline     *string        `json:"headline,omitempty" example:"Full Stack Developer"`
	Bio          *string        `json:"bio,omitempty" example:"Passionate developer with 5+ years of experience"`
	BirthDate    *time.Time     `json:"birth_date,omitempty" example:"1990-01-15T00:00:00Z"`
	Country      *enums.Country `json:"country,omitempty" example:"US"`
	GithubURL    *string        `json:"github_url,omitempty" example:"https://github.com/johndoe"`
	LinkedinURL  *string        `json:"linkedin_url,omitempty" example:"https://linkedin.com/in/johndoe"`
	PitchVideo   *string        `json:"pitch_video,omitempty" example:"https://youtube.com/watch?v=example"`
	Seniority    *string        `json:"seniority,omitempty" example:"Senior"`
	EnglishLevel *string        `json:"english_level,omitempty" example:"Advanced"`
	EFSetScore   *string        `json:"efset_score,omitempty" example:"75"`
	Availability *string        `json:"availability,omitempty" example:"Full-time"`
	Interests    *string        `json:"interests,omitempty" example:"JavaScript, Go, React"`
} // @name UpdateProfileDTO
