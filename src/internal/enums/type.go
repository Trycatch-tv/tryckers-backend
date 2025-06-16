package enums

// Type represents the type of post
//go:generate stringer -type=Type
type Type uint8

const (
	// RegularPost represents a regular post
	RegularPost Type = iota
	// StoryPost represents a story post
	StoryPost
	// VideoPost represents a video post
	VideoPost
)
