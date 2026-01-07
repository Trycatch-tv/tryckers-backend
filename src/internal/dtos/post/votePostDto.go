package dtos

type VotePostDto struct {
	Vote int `json:"vote"` // 1 for upvote, -1 for downvote
}
