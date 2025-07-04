package handlers

import (
	"net/http"

	dto "github.com/Trycatch-tv/tryckers-backend/src/internal/dtos/comment"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentHandler struct {
	Service *services.CommentService
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	var comment dto.CreateCommentDto
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdComment)
}

func (h *CommentHandler) GetCommentsByPostId(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	comments, err := h.Service.GetCommentsByPostId(uuid.Must(uuid.Parse(id)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comments)
}
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	var comment dto.UpdateCommentDto
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	modelComment := models.Comment{
		ID:      uuid.Must(uuid.Parse(id)),
		Content: comment.Content,
	}

	updatedComment, err := h.Service.UpdateComment(&modelComment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedComment)
}
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	_, err := h.Service.DeleteComment(uuid.Must(uuid.Parse(id)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Comment deleted successfully"})
}
