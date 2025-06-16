package services

import (
	"errors"
	"time"

	dt "github.com/Trycatch-tv/tryckers-backend/src/internal/dtos/comment"
	models "github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	repository "github.com/Trycatch-tv/tryckers-backend/src/internal/repository"
	uuid "github.com/google/uuid"
)

type CommentService struct {
	Repo *repository.CommentRepository
}

func (s *CommentService) CreateComment(comment *dt.CreateCommentDto) (models.Comment, error) {
	newComment := models.Comment{
		Content: comment.Content,
		Image:   comment.Image,
		Status:  comment.Status,
		UserID:  comment.UserId,
		PostID:  comment.PostId,
	}
	return s.Repo.CreateComment(&newComment)
}
func (s *CommentService) GetAllComments() ([]dt.ResponseCommentDto, error) {
	var comments []models.Comment
	err := s.Repo.DB.Find(&comments).Error
	commentsDto := make([]dt.ResponseCommentDto, len(comments))
	for i := range comments {
		commentsDto[i] = dt.ResponseCommentDto{
			ID:        comments[i].ID.String(),
			Content:   comments[i].Content,
			Status:    comments[i].Status,
			CreatedAt: comments[i].CreatedAt,
			UpdatedAt: comments[i].UpdatedAt,
			UserId:    comments[i].UserID.String(),
			PostId:    comments[i].PostID.String(),
		}
	}
	return commentsDto, err
}
func (s *CommentService) GetCommentById(id uuid.UUID) (dt.ResponseCommentDto, error) {
	var comment models.Comment
	err := s.Repo.DB.First(&comment, id).Error
	return dt.ResponseCommentDto{
		ID:        comment.ID.String(),
		Content:   comment.Content,
		Status:    comment.Status,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		UserId:    comment.UserID.String(),
		PostId:    comment.PostID.String(),
	}, err
}
func (s *CommentService) UpdateComment(comment *dt.UpdateCommentDto) (models.Comment, error) {
	var updatedComment models.Comment
	commentId := uuid.Must(uuid.Parse(comment.ID))
	commentData, err := s.Repo.GetCommentById(commentId)
	if commentData.ID == uuid.Nil {
		return models.Comment{}, errors.New("Invalid ID")
	}
	if err != nil {
		return models.Comment{}, err
	}
	updatedComment.Content = comment.Content
	updatedComment.Status = comment.Status
	updatedComment.UpdatedAt = time.Now()
	return s.Repo.UpdateComment(&updatedComment)
}
func (s *CommentService) DeleteComment(id uuid.UUID) error {
	err := s.Repo.DeleteComment(id)
	return err
}
