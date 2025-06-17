package services

import (
	"errors"
	"time"

	dtos "github.com/Trycatch-tv/tryckers-backend/src/internal/dtos/post"
	enums "github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	models "github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	repository "github.com/Trycatch-tv/tryckers-backend/src/internal/repository"
	uuid "github.com/google/uuid"
)

type PostService struct {
	Repo *repository.PostRepository
}

func (s *PostService) CreatePost(post *dtos.CreatePostDto) (models.Post, error) {

	newPost := models.Post{
		Title:   post.Title,
		Content: post.Content,
		Image:   post.Image,
		Type:    string(post.Type),
		Tags:    post.Tags,
		Status:  post.Status,
		UserID:  post.UserId,
	}
	return s.Repo.CreatePost(&newPost)
}
func (s *PostService) GetAllPosts() ([]models.Post, error) {
	return s.Repo.GetAllPosts()
}
func (s *PostService) GetPostById(id uuid.UUID) (dtos.ResponsePostDto, error) {
	var post models.Post
	post, err := s.Repo.GetPostById(id)
	return dtos.ResponsePostDto{
		ID:        post.ID.String(),
		Title:     post.Title,
		Content:   post.Content,
		Image:     post.Image,
		Type:      string(post.Type),
		Tags:      post.Tags,
		Status:    enums.PostStatus(string(post.Status)),
		CreatedAt: *post.CreatedAt,
		UpdatedAt: *post.UpdatedAt,
		UserId:    post.UserID.String(),
	}, err
}
func (s *PostService) UpdatePost(post *dtos.UpdatePostDto) (models.Post, error) {
	var updatedPost models.Post
	id := uuid.Must(uuid.Parse(post.ID))
	updatedPost, err := s.Repo.GetPostById(id)
	if id == uuid.Nil {
		return models.Post{}, errors.New("Invalid ID")
	}
	if err != nil {
		return models.Post{}, err
	}
	updatedPost.Title = post.Title
	updatedPost.Content = post.Content
	updatedPost.Image = post.Image
	updatedPost.Type = string(post.Type)
	updatedPost.Tags = post.Tags
	// updatedPost.Status = enums.PostStatus(string(post.Status))
	updatedPost.UserID = uuid.Must(uuid.Parse(post.UserId))
	updatedPost.UpdatedAt = &time.Time{}
	return s.Repo.UpdatePost(&updatedPost)
}
func (s *PostService) DeletePost(id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("Invalid ID")
	}
	return s.Repo.DeletePost(id)
}
