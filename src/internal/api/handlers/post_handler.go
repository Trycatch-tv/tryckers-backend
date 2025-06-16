package handlers

import (
	"net/http"

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

	// Parsea el JSON del body
	if err := c.ShouldBindJSON(&postDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Verifica que UserId no sea cero
	if postDto.UserId == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Se requiere un UserID válido"})
		return
	}

	// Llama al servicio
	createdPost, err := h.Service.CreatePost(&postDto)
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

	post, err := h.Service.GetPostById(uuid.Must(uuid.Parse(id)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
}
func (h *PostHandler) UpdatePost(c *gin.Context) {
	var post dtos.UpdatePostDto

	// Parsea el JSON del body
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	updatedPost, err := h.Service.UpdatePost(&post)
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

	err := h.Service.DeletePost(uuid.Must(uuid.Parse(id)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
