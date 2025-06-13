package enums

type Status string

const (
	PUBLISHED Status = "published"
	EDITED    Status = "edited"
	DRAFT     Status = "draft"
	DELETED   Status = "deleted"
)
