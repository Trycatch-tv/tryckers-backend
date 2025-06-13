package repository

import (
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository struct {
	DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (r *PostRepository) CreatePost(post *models.Post) (models.Post, error) {
	result := r.DB.Create(post)
	return *post, result.Error
}
func (r *PostRepository) GetAllPosts() ([]models.Post, error) {
	var posts []models.Post
	err := r.DB.Find(&posts).Error
	return posts, err
}
func (r *PostRepository) GetPostById(id uuid.UUID) (models.Post, error) {
	var post models.Post
	err := r.DB.First(&post, id).Error
	return post, err
}
func (r *PostRepository) UpdatePost(post *models.Post) (models.Post, error) {

	result := r.DB.Save(post)
	return *post, result.Error
}
func (r *PostRepository) DeletePost(id uuid.UUID) error {
	err := r.DB.Delete(&models.Post{}, id).Error
	return err
}
