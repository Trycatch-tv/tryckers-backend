package enums

type PostStatus string

const (
	PUBLISHED PostStatus = "published"
	EDITED    PostStatus = "edited"
	DRAFT     PostStatus = "draft"
	DELETED   PostStatus = "deleted"
)
