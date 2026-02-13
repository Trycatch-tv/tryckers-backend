package handlers

import (
	"net/http"

	dto "github.com/Trycatch-tv/tryckers-backend/src/internal/dtos/comment"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/services"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	Service *services.CommentService
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	var comment dto.CreateCommentDto
	if err := c.ShouldBindJSON(&comment); err != nil {
		HandleBadRequest(c, "datos de entrada inválidos")
		return
	}
	modelComment := models.Comment{
		Content: comment.Content,
		Image:   comment.Image,
		PostID:  comment.PostId,
		UserID:  comment.UserId,
		Status:  bool(enums.Active),
	}
	createdComment, err := h.Service.CreateComment(&modelComment)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, createdComment)
}

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
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	var comment dto.UpdateCommentDto
	if err := c.ShouldBindJSON(&comment); err != nil {
		HandleBadRequest(c, "datos de entrada inválidos")
		return
	}
	id := c.Param("id")
	parsedId, err := ParseUUID(id)
	if err != nil {
		HandleError(c, err)
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
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id := c.Param("id")

	parsedId, err := ParseUUID(id)
	if err != nil {
		HandleError(c, err)
		return
	}

	_, err = h.Service.DeleteComment(parsedId)
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Comment deleted successfully"})
}
