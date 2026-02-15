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
		HandleBindingError(c, err)
		return
	}

	if postDto.UserID == uuid.Nil {
		HandleBadRequest(c, "se requiere un UserID válido")
		return
	}

	// Validaciones adicionales del DTO
	if err := postDto.Validate(); err != nil {
		HandleBadRequest(c, err.Error())
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
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dtos.ToResponsePostDto(&createdPost, 0))
}
func (h *PostHandler) GetAllPosts(c *gin.Context) {

	posts, err := h.Service.GetAllPosts()
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, posts)
}
func (h *PostHandler) GetPostById(c *gin.Context) {
	id := c.Param("id")

	parsedId, err := ParseUUID(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	var loggedUserId *uuid.UUID
	if userIdStr, exists := c.Get("userId"); exists {
		if uidStr, ok := userIdStr.(string); ok {
			uid, err := uuid.Parse(uidStr)
			if err == nil {
				loggedUserId = &uid
			}
		}
	}

	post, userVote, err := h.Service.GetPostById(parsedId, loggedUserId)
	if err != nil {
		HandleError(c, err)
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
		HandleBindingError(c, err)
		return
	}

	// Verificar que hay cambios
	if !post.HasChanges() {
		HandleBadRequest(c, dtos.ErrNoChangesProvided.Error())
		return
	}

	// Validaciones adicionales
	if err := post.Validate(); err != nil {
		HandleBadRequest(c, err.Error())
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
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, dtos.ToResponsePostDto(&result, 0))
}
func (h *PostHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")

	parsedId, err := ParseUUID(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	_, err = h.Service.DeletePost(parsedId)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Post deleted successfully"})
}

func (h *PostHandler) GetPostsByUserId(c *gin.Context) {
	ownerId := c.Param("ownerId")
	parsedOwnerId, err := ParseUUID(ownerId)
	if err != nil {
		HandleError(c, err)
		return
	}
	var loggedUserId *uuid.UUID
	if v, exists := c.Get("userId"); exists {
		if uidStr, ok := v.(string); ok {
			id, err := uuid.Parse(uidStr)
			if err == nil {
				loggedUserId = &id
			}
		}
	}
	posts, userVotes, err := h.Service.GetPostsByUserId(parsedOwnerId, loggedUserId)
	if err != nil {
		HandleError(c, err)
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
		HandleUnauthorized(c, "usuario no autenticado")
		return
	}
	var voteDto dtos.VotePostDto
	if err := c.ShouldBindJSON(&voteDto); err != nil {
		HandleBindingError(c, err)
		return
	}
	parsedPostId, err := ParseUUID(postId)
	if err != nil {
		HandleError(c, err)
		return
	}
	parsedUserId, err := ParseUUID(userId.(string))
	if err != nil {
		HandleError(c, err)
		return
	}

	postVote, err := h.Service.PostVote(parsedPostId, parsedUserId, int8(voteDto.Vote))
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, postVote)
}
