package services

import (
	"errors"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	models "github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	repository "github.com/Trycatch-tv/tryckers-backend/src/internal/repository"
	uuid "github.com/google/uuid"
)

type PostService struct {
	Repo *repository.PostRepository
}

func (s *PostService) CreatePost(post models.Post) (models.Post, error) {
	return s.Repo.CreatePost(post)
}
func (s *PostService) GetAllPosts() ([]models.Post, error) {
	return s.Repo.GetAllPosts()
}
func (s *PostService) GetPostById(id uuid.UUID) (models.Post, error) {
	return s.Repo.GetPostById(id)
}
func (s *PostService) UpdatePost(post models.Post) (models.Post, error) {
	var updatedPost models.Post

	updatedPost, err := s.Repo.GetPostById(post.ID)
	if err != nil {
		return models.Post{}, err
	}
	updatedPost.Title = post.Title
	updatedPost.Content = post.Content
	updatedPost.Image = post.Image
	updatedPost.Type = post.Type
	updatedPost.Tags = post.Tags
	updatedPost.Status = enums.PostStatus(string(post.Status))
	return s.Repo.UpdatePost(&updatedPost)
}
func (s *PostService) DeletePost(id uuid.UUID) (models.Post, error) {
	if id == uuid.Nil {
		return models.Post{}, errors.New("invalid ID")
	}
	post, err := s.Repo.GetPostById(id)
	if err != nil {
		return models.Post{}, err
	}

	post.Status = enums.DELETED
	return s.Repo.DeletePost(&post)
}
