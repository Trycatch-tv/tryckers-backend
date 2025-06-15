package services

import (
	"time"

	dt "github.com/Trycatch-tv/tryckers-backend/src/internal/dtos"
	enums "github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	models "github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	repository "github.com/Trycatch-tv/tryckers-backend/src/internal/repository"
	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentService struct {
	DB *gorm.DB
}

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{DB: db}
}

func (s *CommentService) CreateComment(comment *dt.CreateCommentRequest) (models.Comment, error) {
	newComment := models.Comment{
		Content: comment.Content,
		// Status:    string(comment.Status),
		UserID: comment.UserId,
		PostID: comment.PostId,
		// CreatedAt: time.Now(),
	}
	return repository.NewCommentRepository(s.DB).CreateComment(&newComment)
}
func (s *CommentService) GetAllComments() ([]dt.CommentDto, error) {
	var comments []models.Comment
	err := s.DB.Find(&comments).Error
	commentsDto := make([]dt.CommentDto, len(comments))
	for i := range comments {
		commentsDto[i] = dt.CommentDto{
			ID:        comments[i].ID.String(),
			Content:   comments[i].Content,
			Status:    enums.Status(comments[i].Status),
			CreatedAt: comments[i].CreatedAt,
			UpdatedAt: comments[i].UpdatedAt,
			UserId:    comments[i].UserID.String(),
			PostId:    comments[i].PostID.String(),
		}
	}
	return commentsDto, err
}
func (s *CommentService) GetCommentById(id uuid.UUID) (dt.CommentDto, error) {
	var comment models.Comment
	err := s.DB.First(&comment, id).Error
	return dt.CommentDto{
		ID:        comment.ID.String(),
		Content:   comment.Content,
		Status:    enums.Status(comment.Status),
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
		UserId:    comment.UserID.String(),
		PostId:    comment.PostID.String(),
	}, err
}
func (s *CommentService) UpdateComment(comment *dt.CreateCommentRequest) (models.Comment, error) {
	var updatedComment models.Comment
	err := s.DB.First(&updatedComment, comment.ID).Error
	if err != nil {
		return models.Comment{}, err
	}
	updatedComment.Content = comment.Content
	updatedComment.Status = string(comment.Status)
	updatedComment.UpdatedAt = time.Now()
	return repository.NewCommentRepository(s.DB).UpdateComment(&updatedComment)
}
func (s *CommentService) DeleteComment(id uuid.UUID) error {
	err := s.DB.Delete(&models.Comment{}, id).Error
	return err
}
