package handlers

import (
	"net/http"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"

	dtos "github.com/Trycatch-tv/tryckers-backend/src/internal/dtos/post"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostHandler struct {
	Service *services.PostService
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var postDto dtos.CreatePostDto

	if err := c.ShouldBindJSON(&postDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if postDto.UserId == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere un UserID válido"})
		return
	}

	newPost := models.Post{
		Title:   postDto.Title,
		Content: postDto.Content,
		Image:   postDto.Image,
		Type:    postDto.Type,
		Tags:    postDto.Tags,
		Status:  postDto.Status,
		UserID:  postDto.UserId,
	}

	createdPost, err := h.Service.CreatePost(newPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdPost)
}
func (h *PostHandler) GetAllPosts(c *gin.Context) {

	posts, err := h.Service.GetAllPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}
func (h *PostHandler) GetPostById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var loggedUserId *uuid.UUID
	if userIdStr, exists := c.Get("userId"); exists {
		if uidStr, ok := userIdStr.(string); ok {
			uid := uuid.Must(uuid.Parse(uidStr))
			loggedUserId = &uid
		}
	}

	post, userVote, err := h.Service.GetPostById(uuid.Must(uuid.Parse(id)), loggedUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := gin.H{
		"id":          post.ID,
		"title":       post.Title,
		"content":     post.Content,
		"image":       post.Image,
		"type":        post.Type,
		"tags":        post.Tags,
		"status":      post.Status,
		"created_at":  post.CreatedAt,
		"updated_at":  post.UpdatedAt,
		"user_id":     post.UserID,
		"user":        post.User,
		"votes_count": post.VotesCount,
		"user_vote":   userVote,
	}
	c.JSON(http.StatusOK, response)
}
func (h *PostHandler) UpdatePost(c *gin.Context) {
	var post dtos.UpdatePostDto

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPost := models.Post{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Image:   post.Image,
		Type:    post.Type,
		Tags:    post.Tags,
		Status:  post.Status,
	}

	updatedPost, err := h.Service.UpdatePost(updatedPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedPost)
}
func (h *PostHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	_, err := h.Service.DeletePost(uuid.Must(uuid.Parse(id)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Post deleted successfully"})
}

func (h *PostHandler) GetPostsByUserId(c *gin.Context) {
	ownerId := c.Param("ownerId")
	if ownerId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Owner ID"})
		return
	}
	var loggedUserId *uuid.UUID
	if v, exists := c.Get("userId"); exists {
		id := uuid.Must(uuid.Parse(v.(string)))
		loggedUserId = &id
	}
	posts, userVotes, err := h.Service.GetPostsByUserId(uuid.Must(uuid.Parse(ownerId)), loggedUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Mapear a ResponsePostDto y agregar el campo user_vote
	var resp []dtos.ResponsePostDto
	for _, p := range posts {
		dto := dtos.ResponsePostDto{
			ID:         p.ID.String(),
			Title:      p.Title,
			Content:    p.Content,
			Image:      p.Image,
			Type:       p.Type,
			Tags:       p.Tags,
			Status:     p.Status,
			CreatedAt:  *p.CreatedAt,
			UpdatedAt:  *p.UpdatedAt,
			UserId:     p.UserID.String(),
			UserVote:   0,
			VotesCount: p.VotesCount,
		}
		if loggedUserId != nil {
			if v, ok := userVotes[p.ID]; ok {
				dto.UserVote = v
			}
		}
		resp = append(resp, dto)
	}
	c.JSON(http.StatusOK, resp)
}

func (h *PostHandler) PostVote(c *gin.Context) {
	postId := c.Param("id")
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
		return
	}
	var voteDto dtos.VotePostDto
	if err := c.ShouldBindJSON(&voteDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if postId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Post ID"})
		return
	}

	postVote, err := h.Service.PostVote(uuid.Must(uuid.Parse(postId)), uuid.Must(uuid.Parse(userId.(string))), int8(voteDto.Vote))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, postVote)
}
