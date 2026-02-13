package services

import (
	"errors"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	apperrors "github.com/Trycatch-tv/tryckers-backend/src/internal/errors"
	models "github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	repository "github.com/Trycatch-tv/tryckers-backend/src/internal/repository"
	uuid "github.com/google/uuid"
	"gorm.io/gorm"
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
func (s *PostService) GetPostById(id uuid.UUID, loggedUserId *uuid.UUID) (models.Post, int8, error) {
	post, vote, err := s.Repo.GetPostById(id, loggedUserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Post{}, 0, apperrors.ErrPostNotFound
		}
		return models.Post{}, 0, apperrors.NewInternalError("error al obtener post", err)
	}
	return post, vote, nil
}
func (s *PostService) UpdatePost(post models.Post) (models.Post, error) {
	var updatedPost models.Post

	updatedPost, _, err := s.Repo.GetPostById(post.ID, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Post{}, apperrors.ErrPostNotFound
		}
		return models.Post{}, apperrors.NewInternalError("error al buscar post", err)
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
		return models.Post{}, apperrors.ErrInvalidID
	}
	post, _, err := s.Repo.GetPostById(id, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Post{}, apperrors.ErrPostNotFound
		}
		return models.Post{}, apperrors.NewInternalError("error al buscar post", err)
	}

	post.Status = enums.DELETED
	return s.Repo.DeletePost(&post)
}

// Devuelve los posts de un usuario y el voto del usuario logueado (si se pasa loggedUserId)
func (s *PostService) GetPostsByUserId(userId uuid.UUID, loggedUserId *uuid.UUID) ([]models.Post, map[uuid.UUID]int8, error) {
	return s.Repo.GetPostsByUserId(userId, loggedUserId)
}

func (s *PostService) PostVote(postId uuid.UUID, userId uuid.UUID, vote int8) (models.PostVote, error) {
	postvote, err := s.Repo.PostVote(postId, userId, vote)
	if err != nil {
		return models.PostVote{}, err
	}
	return postvote, nil
}
