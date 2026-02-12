package handlers

import (
	"net/http"
	"strings"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"

	dtos "github.com/Trycatch-tv/tryckers-backend/src/internal/dtos/post"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostHandler struct {
	Service *services.PostService
}

// tagsToString convierte un slice de tags a una cadena separada por comas
func tagsToString(tags []string) string {
	return strings.Join(tags, ",")
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var postDto dtos.CreatePostDto

	if err := c.ShouldBindJSON(&postDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if postDto.UserID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere un UserID válido"})
		return
	}

	// Validaciones adicionales del DTO
	if err := postDto.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newPost := models.Post{
		Title:    postDto.Title,
		Content:  postDto.Content,
		Image:    postDto.Image,
		Type:     postDto.Type,
		Tags:     tagsToString(postDto.Tags),
		Status:   postDto.Status,
		UserID:   postDto.UserID,
		MediaURL: postDto.MediaURL,
	}

	createdPost, err := h.Service.CreatePost(newPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dtos.ToResponsePostDto(&createdPost, 0))
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

	// Verificar que hay cambios
	if !post.HasChanges() {
		c.JSON(http.StatusBadRequest, gin.H{"error": dtos.ErrNoChangesProvided.Error()})
		return
	}

	// Validaciones adicionales
	if err := post.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Construir el modelo con solo los campos proporcionados
	updatedPost := models.Post{
		ID: post.ID,
	}

	if post.Title != nil {
		updatedPost.Title = *post.Title
	}
	if post.Content != nil {
		updatedPost.Content = *post.Content
	}
	if post.Image != nil {
		updatedPost.Image = *post.Image
	}
	if post.Type != nil {
		updatedPost.Type = *post.Type
	}
	if post.Tags != nil {
		updatedPost.Tags = tagsToString(post.Tags)
	}
	if post.Status != nil {
		updatedPost.Status = *post.Status
	}
	if post.MediaURL != nil {
		updatedPost.MediaURL = *post.MediaURL
	}

	result, err := h.Service.UpdatePost(updatedPost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dtos.ToResponsePostDto(&result, 0))
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

	// Convertir userVotes de map[uuid.UUID]int8 a map[string]int8
	userVotesMap := make(map[string]int8)
	for id, vote := range userVotes {
		userVotesMap[id.String()] = vote
	}

	// Usar la función de conversión del DTO
	response := dtos.ToResponsePostListDto(posts, int64(len(posts)), 1, len(posts), userVotesMap)
	c.JSON(http.StatusOK, response.Posts)
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
