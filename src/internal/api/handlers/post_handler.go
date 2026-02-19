package handlers

import (
	"net/http"
	"strings"

	apperrors "github.com/Trycatch-tv/tryckers-backend/src/internal/errors"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"

	dtos "github.com/Trycatch-tv/tryckers-backend/src/internal/dtos/post"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
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

// getAuthUser extrae userId y role del contexto del middleware de autenticación
func getAuthUser(c *gin.Context) (uuid.UUID, enums.UserRole, error) {
	userIdStr, exists := c.Get("userId")
	if !exists {
		return uuid.Nil, "", apperrors.ErrUnauthorized
	}
	userId, err := uuid.Parse(userIdStr.(string))
	if err != nil {
		return uuid.Nil, "", apperrors.ErrUnauthorized
	}
	roleStr, _ := c.Get("role")
	role := enums.UserRole(roleStr.(string))
	return userId, role, nil
}

// CreatePost godoc
// @Summary      Create a new post
// @Description  Create a new post. The user_id is taken from the JWT token.
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Param        post  body      dtos.CreatePostDto  true  "Post data"
// @Success      201  {object}  dtos.ResponsePostDto  "Created post"
// @Failure      400  {object}  ErrorResponse  "Invalid input"
// @Failure      401  {object}  ErrorResponse  "Unauthorized"
// @Security     BearerAuth
// @Router       /posts [post]
func (h *PostHandler) CreatePost(c *gin.Context) {
	var postDto dtos.CreatePostDto

	if err := c.ShouldBindJSON(&postDto); err != nil {
		HandleBindingError(c, err)
		return
	}

	// Obtener user_id del token JWT (no del body)
	userId, _, err := getAuthUser(c)
	if err != nil {
		HandleError(c, err)
		return
	}
	postDto.UserID = userId

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
// GetAllPosts godoc
// @Summary      Get all posts
// @Description  Retrieve a list of all posts (excluding deleted)
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Post  "List of posts"
// @Failure      500  {object}  ErrorResponse  "Internal server error"
// @Security     BearerAuth
// @Router       /posts [get]
func (h *PostHandler) GetAllPosts(c *gin.Context) {

	posts, err := h.Service.GetAllPosts()
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, posts)
}
// GetPostById godoc
// @Summary      Get a post by ID
// @Description  Retrieve a specific post by its ID, including the logged user's vote
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Post ID (UUID)"
// @Success      200  {object}  dtos.ResponsePostDto  "Post details"
// @Failure      400  {object}  ErrorResponse  "Invalid ID"
// @Failure      404  {object}  ErrorResponse  "Post not found"
// @Security     BearerAuth
// @Router       /posts/{id} [get]
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
// UpdatePost godoc
// @Summary      Update a post
// @Description  Update an existing post. Only the post author or an admin can update it.
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Param        post  body      dtos.UpdatePostDto  true  "Post update data"
// @Success      200  {object}  dtos.ResponsePostDto  "Updated post"
// @Failure      400  {object}  ErrorResponse  "Invalid input"
// @Failure      403  {object}  ErrorResponse  "Forbidden"
// @Failure      404  {object}  ErrorResponse  "Post not found"
// @Security     BearerAuth
// @Router       /posts [put]
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

	// Verificar propiedad: solo el autor o un admin pueden actualizar
	userId, role, err := getAuthUser(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	existingPost, _, err := h.Service.GetPostById(post.ID, nil)
	if err != nil {
		HandleError(c, err)
		return
	}

	if existingPost.UserID != userId && role != enums.Admin {
		HandleError(c, apperrors.ErrForbidden)
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
// DeletePost godoc
// @Summary      Delete a post
// @Description  Soft-delete a post (sets status to deleted). Only the post author or an admin can delete it.
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Post ID (UUID)"
// @Success      204  "Post deleted successfully"
// @Failure      400  {object}  ErrorResponse  "Invalid ID"
// @Failure      403  {object}  ErrorResponse  "Forbidden"
// @Failure      404  {object}  ErrorResponse  "Post not found"
// @Security     BearerAuth
// @Router       /posts/{id} [delete]
func (h *PostHandler) DeletePost(c *gin.Context) {
	id := c.Param("id")

	parsedId, err := ParseUUID(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Verificar propiedad: solo el autor o un admin pueden eliminar
	userId, role, err := getAuthUser(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	existingPost, _, err := h.Service.GetPostById(parsedId, nil)
	if err != nil {
		HandleError(c, err)
		return
	}

	if existingPost.UserID != userId && role != enums.Admin {
		HandleError(c, apperrors.ErrForbidden)
		return
	}

	_, err = h.Service.DeletePost(parsedId)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Post deleted successfully"})
}

// GetPostsByUserId godoc
// @Summary      Get posts by user ID
// @Description  Retrieve all posts from a specific user
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Param        ownerId  path      string  true  "User ID (UUID)"
// @Success      200  {array}   dtos.ResponsePostDto  "List of user posts"
// @Failure      400  {object}  ErrorResponse  "Invalid ID"
// @Failure      500  {object}  ErrorResponse  "Internal server error"
// @Security     BearerAuth
// @Router       /users/{ownerId}/posts [get]
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

// PostVote godoc
// @Summary      Vote on a post
// @Description  Upvote or downvote a post. Vote value: 1 for upvote, -1 for downvote, 0 to remove vote.
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Param        id    path      string          true  "Post ID (UUID)"
// @Param        vote  body      dtos.VotePostDto  true  "Vote data"
// @Success      200  {object}  models.PostVote  "Vote result"
// @Failure      400  {object}  ErrorResponse  "Invalid input"
// @Failure      401  {object}  ErrorResponse  "Unauthorized"
// @Security     BearerAuth
// @Router       /posts/{id}/vote [post]
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

// Cartelera godoc
// @Summary      Get top posts of the week
// @Description  Retrieve the most popular posts of the last 7 days, ordered by votes
// @Tags         Posts
// @Accept       json
// @Produce      json
// @Success      200  {array}   dtos.ResponsePostDto  "Top posts of the week"
// @Failure      500  {object}  ErrorResponse  "Internal server error"
// @Security     BearerAuth
// @Router       /cartelera [get]
func (h *PostHandler) Cartelera(c *gin.Context) {
	posts, err := h.Service.GetCartelera(10)
	if err != nil {
		HandleError(c, err)
		return
	}

	emptyVotes := make(map[string]int8)
	response := dtos.ToResponsePostListDto(posts, int64(len(posts)), 1, len(posts), emptyVotes)
	c.JSON(http.StatusOK, response.Posts)
}
