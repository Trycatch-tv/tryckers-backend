package repository

import (
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentRepository struct {
	DB *gorm.DB
}

func (r *CommentRepository) CreateComment(comment *models.Comment) (models.Comment, error) {
	result := r.DB.Create(comment)
	if result.Error != nil {
		return models.Comment{}, result.Error
	}

	var createdComment models.Comment
	err := r.DB.Preload("User").Preload("Post").Preload("Post.User").First(&createdComment, comment.ID).Error
	return createdComment, err
}
func (r *CommentRepository) GetAllComments() ([]models.Comment, error) {
	var comments []models.Comment
	err := r.DB.Where("status != ?", enums.Inactive).Find(&comments).Error
	return comments, err
}
func (r *CommentRepository) GetCommentById(id uuid.UUID) (models.Comment, error) {
	var comment models.Comment
	err := r.DB.First(&comment, id).Error
	return comment, err
}
func (r *CommentRepository) UpdateComment(comment *models.Comment) (models.Comment, error) {
	result := r.DB.Save(comment)
	return *comment, result.Error
}
func (r *CommentRepository) DeleteComment(id uuid.UUID) error {
	err := r.DB.Save(&models.Comment{ID: id, Status: bool(enums.Inactive)})
	return err.Error
}
