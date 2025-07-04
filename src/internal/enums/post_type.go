package enums

// Type represents the type of post
//
//go:generate stringer -type=PostType
type PostType string

const (
	RegularPost PostType = "regular"
	StoryPost   PostType = "story"
	VideoPost   PostType = "video"
)
