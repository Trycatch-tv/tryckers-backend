package services

import (
	"errors"

	models "github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	repository "github.com/Trycatch-tv/tryckers-backend/src/internal/repository"
	uuid "github.com/google/uuid"
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
		return []models.Comment{}, err
	}
	return comments, nil
}
func (s *CommentService) UpdateComment(comment *models.Comment) (models.Comment, error) {
	updatedComment, err := s.Repo.GetCommentById(comment.ID)
	if updatedComment.ID == uuid.Nil {
		return models.Comment{}, errors.New("Invalid ID")
	}
	if err != nil {
		return models.Comment{}, err
	}
	updatedComment.Content = comment.Content
	updatedComment.Status = comment.Status
	updatedComment.Image = comment.Image
	return s.Repo.UpdateComment(&updatedComment)
}

func (s *CommentService) DeleteComment(id uuid.UUID) error {
	err := s.Repo.DeleteComment(id)
	return err
}
