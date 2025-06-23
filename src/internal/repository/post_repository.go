package repository

import (
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository struct {
	DB *gorm.DB
}

func (r *PostRepository) CreatePost(post models.Post) (models.Post, error) {
	result := r.DB.Create(post)
	if result.Error != nil {
		return post, result.Error
	}
	var createdPost models.Post
	err := r.DB.Preload("User").First(&createdPost, post.ID).Error
	return createdPost, err
}
func (r *PostRepository) GetAllPosts() ([]models.Post, error) {
	var posts []models.Post
	err := r.DB.Preload("User").Where("status != ?", enums.DELETED).Find(&posts).Error
	return posts, err
}
func (r *PostRepository) GetPostById(id uuid.UUID) (models.Post, error) {
	var post models.Post
	err := r.DB.Where("status != ?", enums.DELETED).First(&post, id).Error
	return post, err
}
func (r *PostRepository) UpdatePost(post *models.Post) (models.Post, error) {
	result := r.DB.Save(post)
	return *post, result.Error
}
func (r *PostRepository) DeletePost(post *models.Post) (models.Post, error) {
	result := r.DB.Save(post)
	return *post, result.Error
}
