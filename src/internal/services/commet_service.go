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
func (s *CommentService) GetAllComments() ([]models.Comment, error) {
	var comments []models.Comment
	err := s.Repo.DB.Find(&comments).Error
	return comments, err
}
func (s *CommentService) GetCommentById(id uuid.UUID) (models.Comment, error) {
	var comment models.Comment
	err := s.Repo.DB.First(&comment, id).Error
	return comment, err
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
