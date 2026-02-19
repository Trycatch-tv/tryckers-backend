package handlers

import (
	"net/http"

	dto "github.com/Trycatch-tv/tryckers-backend/src/internal/dtos/comment"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	apperrors "github.com/Trycatch-tv/tryckers-backend/src/internal/errors"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/services"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	Service     *services.CommentService
	PostService *services.PostService
}

// CreateComment godoc
// @Summary      Create a new comment
// @Description  Create a new comment on a post. The user_id is taken from the JWT token.
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Param        comment  body      dto.CreateCommentDto  true  "Comment data"
// @Success      201  {object}  models.Comment  "Created comment"
// @Failure      400  {object}  ErrorResponse  "Invalid input"
// @Failure      401  {object}  ErrorResponse  "Unauthorized"
// @Security     BearerAuth
// @Router       /comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	var comment dto.CreateCommentDto
	if err := c.ShouldBindJSON(&comment); err != nil {
		HandleBindingError(c, err)
		return
	}

	// Obtener user_id del token JWT (no del body)
	userId, _, err := getAuthUser(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	modelComment := models.Comment{
		Content: comment.Content,
		Image:   comment.Image,
		PostID:  comment.PostId,
		UserID:  userId,
		Status:  bool(enums.Active),
	}
	createdComment, err := h.Service.CreateComment(&modelComment)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdComment)
}

// GetCommentsByPostId godoc
// @Summary      Get comments by post ID
// @Description  Retrieve all comments for a specific post
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Post ID (UUID)"
// @Success      200  {array}   models.Comment  "List of comments"
// @Failure      400  {object}  ErrorResponse  "Invalid ID"
// @Failure      500  {object}  ErrorResponse  "Internal server error"
// @Security     BearerAuth
// @Router       /posts/{id}/comments [get]
func (h *CommentHandler) GetCommentsByPostId(c *gin.Context) {
	id := c.Param("id")

	parsedId, err := ParseUUID(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	comments, err := h.Service.GetCommentsByPostId(parsedId)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, comments)
}

// UpdateComment godoc
// @Summary      Update a comment
// @Description  Update an existing comment. Only the comment author or an admin can update it.
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Param        id       path      string                true  "Comment ID (UUID)"
// @Param        comment  body      dto.UpdateCommentDto  true  "Comment update data"
// @Success      200  {object}  models.Comment  "Updated comment"
// @Failure      400  {object}  ErrorResponse  "Invalid input"
// @Failure      403  {object}  ErrorResponse  "Forbidden"
// @Failure      404  {object}  ErrorResponse  "Comment not found"
// @Security     BearerAuth
// @Router       /comments/{id} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	var comment dto.UpdateCommentDto
	if err := c.ShouldBindJSON(&comment); err != nil {
		HandleBindingError(c, err)
		return
	}
	id := c.Param("id")
	parsedId, err := ParseUUID(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Verificar propiedad: solo el autor del comentario o un admin pueden actualizar
	userId, role, err := getAuthUser(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	existingComment, err := h.Service.GetCommentById(parsedId)
	if err != nil {
		HandleError(c, err)
		return
	}

	if existingComment.UserID != userId && role != enums.Admin {
		HandleError(c, apperrors.ErrForbidden)
		return
	}

	modelComment := models.Comment{
		ID:      parsedId,
		Content: comment.Content,
	}

	updatedComment, err := h.Service.UpdateComment(&modelComment)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, updatedComment)
}

// DeleteComment godoc
// @Summary      Delete a comment
// @Description  Soft-delete a comment. The comment author, post owner, or an admin can delete it.
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Comment ID (UUID)"
// @Success      204  "Comment deleted successfully"
// @Failure      400  {object}  ErrorResponse  "Invalid ID"
// @Failure      403  {object}  ErrorResponse  "Forbidden"
// @Failure      404  {object}  ErrorResponse  "Comment not found"
// @Security     BearerAuth
// @Router       /comments/{id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id := c.Param("id")

	parsedId, err := ParseUUID(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	// Verificar propiedad: autor del comentario, dueño del post, o admin
	userId, role, err := getAuthUser(c)
	if err != nil {
		HandleError(c, err)
		return
	}

	existingComment, err := h.Service.GetCommentById(parsedId)
	if err != nil {
		HandleError(c, err)
		return
	}

	isCommentAuthor := existingComment.UserID == userId
	isAdmin := role == enums.Admin

	// Verificar si es dueño del post
	isPostOwner := false
	if existingComment.Post != nil {
		isPostOwner = existingComment.Post.UserID == userId
	} else {
		// Si no se precargó el post, buscarlo
		post, _, postErr := h.PostService.GetPostById(existingComment.PostID, nil)
		if postErr == nil {
			isPostOwner = post.UserID == userId
		}
	}

	if !isCommentAuthor && !isPostOwner && !isAdmin {
		HandleError(c, apperrors.ErrForbidden)
		return
	}

	_, err = h.Service.DeleteComment(parsedId)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Comment deleted successfully"})
}
