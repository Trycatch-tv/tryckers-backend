package services

import (
	"errors"

	apperrors "github.com/Trycatch-tv/tryckers-backend/src/internal/errors"
	models "github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	repository "github.com/Trycatch-tv/tryckers-backend/src/internal/repository"
	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentService struct {
	Repo *repository.CommentRepository
}

func (s *CommentService) CreateComment(comment *models.Comment) (models.Comment, error) {
	return s.Repo.CreateComment(comment)
}

func (s *CommentService) GetCommentsByPostId(id uuid.UUID) ([]models.Comment, error) {
	comments, err := s.Repo.GetCommentsByPostId(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []models.Comment{}, nil
		}
		return []models.Comment{}, apperrors.NewInternalError("error al obtener comentarios", err)
	}
	return comments, nil
}
func (s *CommentService) UpdateComment(comment *models.Comment) (models.Comment, error) {
	updatedComment, err := s.Repo.GetCommentById(comment.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Comment{}, apperrors.ErrCommentNotFound
		}
		return models.Comment{}, apperrors.NewInternalError("error al buscar comentario", err)
	}
	if updatedComment.ID == uuid.Nil {
		return models.Comment{}, apperrors.ErrCommentNotFound
	}
	updatedComment.Content = comment.Content
	return s.Repo.UpdateComment(&updatedComment)
}

func (s *CommentService) DeleteComment(id uuid.UUID) (models.Comment, error) {
	if id == uuid.Nil {
		return models.Comment{}, apperrors.ErrInvalidID
	}
	comment, err := s.Repo.GetCommentById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Comment{}, apperrors.ErrCommentNotFound
		}
		return models.Comment{}, apperrors.NewInternalError("error al buscar comentario", err)
	}

	return s.Repo.DeleteComment(comment.ID)
}
