package repository

import (
	"fmt"

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
		return models.Comment{}, fmt.Errorf("error al crear comentario: %w", result.Error)
	}

	var createdComment models.Comment
	err := r.DB.Preload("User").Preload("Post").Preload("Post.User").First(&createdComment, comment.ID).Error
	if err != nil {
		return models.Comment{}, fmt.Errorf("error al obtener comentario creado: %w", err)
	}
	return createdComment, nil
}

func (r *CommentRepository) GetCommentById(id uuid.UUID) (models.Comment, error) {
	var comment models.Comment
	err := r.DB.Preload("User").Preload("Post").Preload("Post.User").First(&comment, id).Error
	if err != nil {
		return models.Comment{}, err
	}
	return comment, nil
}

func (r *CommentRepository) GetCommentsByPostId(id uuid.UUID) ([]models.Comment, error) {
	var comments []models.Comment
	err := r.DB.Preload("User").Where("post_id", id).Find(&comments).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener comentarios del post: %w", err)
	}
	return comments, nil
}

func (r *CommentRepository) UpdateComment(comment *models.Comment) (models.Comment, error) {
	result := r.DB.Save(comment)
	if result.Error != nil {
		return models.Comment{}, fmt.Errorf("error al actualizar comentario: %w", result.Error)
	}
	return *comment, nil
}

func (r *CommentRepository) DeleteComment(id uuid.UUID) (models.Comment, error) {
	comment := models.Comment{ID: id, Status: bool(enums.Inactive)}
	result := r.DB.Save(&comment)
	if result.Error != nil {
		return models.Comment{}, fmt.Errorf("error al eliminar comentario: %w", result.Error)
	}
	return comment, nil
}
