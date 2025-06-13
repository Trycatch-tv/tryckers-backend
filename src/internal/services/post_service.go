package services

import (
	"time"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/dtos"
	d "github.com/Trycatch-tv/tryckers-backend/src/internal/dtos"
	enums "github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	models "github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	repository "github.com/Trycatch-tv/tryckers-backend/src/internal/repository"
	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type PostService struct {
	DB *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{DB: db}
}

func (s *PostService) CreatePost(post *dtos.CreatePostRequest) (models.Post, error) {
	newPost := models.Post{
		Title:     post.Title,
		Content:   post.Content,
		Status:    string(post.Status),
		UserID:    post.UserId,
		CreatedAt: &time.Time{},
	}
	return repository.NewPostRepository(s.DB).CreatePost(&newPost)
}
func (s *PostService) GetAllPosts() ([]d.PostDto, error) {
	var posts []models.Post
	err := s.DB.Find(&posts).Error
	if err != nil {
		return nil, err
	}

	dtos := make([]d.PostDto, len(posts))
	for i := range posts {
		dtos[i] = d.PostDto{
			ID:        posts[i].ID.String(),
			Title:     posts[i].Title,
			Content:   posts[i].Content,
			Status:    enums.Status(posts[i].Status),
			CreatedAt: *posts[i].CreatedAt,
			UpdatedAt: *posts[i].UpdatedAt,
			UserId:    posts[i].UserID.String(),
		}
	}
	return dtos, nil
}
func (s *PostService) GetPostById(id uuid.UUID) (d.PostDto, error) {
	var post models.Post
	err := s.DB.First(&post, id).Error
	return dtos.PostDto{
		ID:        post.ID.String(),
		Title:     post.Title,
		Content:   post.Content,
		Status:    enums.Status(post.Status),
		CreatedAt: *post.CreatedAt,
		UpdatedAt: *post.UpdatedAt,
		UserId:    post.UserID.String(),
	}, err
}
func (s *PostService) UpdatePost(post *dtos.CreatePostRequest) (models.Post, error) {
	var updatedPost models.Post
	err := s.DB.First(&updatedPost, post.ID).Error
	if err != nil {
		return models.Post{}, err
	}
	updatedPost.Title = post.Title
	updatedPost.Content = post.Content
	updatedPost.Status = string(post.Status)
	updatedPost.UserID = post.UserId
	updatedPost.UpdatedAt = &time.Time{}
	return repository.NewPostRepository(s.DB).UpdatePost(&updatedPost)
}
func (s *PostService) DeletePost(id uuid.UUID) error {
	return repository.NewPostRepository(s.DB).DeletePost(id)
}
