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
func (h *CommentHandler) GetAllComments(c *gin.Context) {
	comments, err := h.Service.GetAllComments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comments)
}
func (h *CommentHandler) GetCommentById(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	comment, err := h.Service.GetCommentById(uuid.Must(uuid.Parse(id)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, comment)
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
		ID:      comment.ID,
		Content: comment.Content,
		Image:   comment.Image,
		Status:  bool(enums.Active),
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

	err := h.Service.DeleteComment(uuid.Must(uuid.Parse(id)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
