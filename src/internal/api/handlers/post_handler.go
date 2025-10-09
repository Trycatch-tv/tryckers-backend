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

	post, err := h.Service.GetPostById(uuid.Must(uuid.Parse(id)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, post)
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
