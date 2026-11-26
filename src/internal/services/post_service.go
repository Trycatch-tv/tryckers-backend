package services

import (
	"context"
	"errors"
	"log"
	"mime/multipart"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	apperrors "github.com/Trycatch-tv/tryckers-backend/src/internal/errors"
	models "github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/services/storage"
	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

const maxPostImageBytes = 5 * 1024 * 1024

type PostRepositoryInterface interface {
	CreatePost(post models.Post) (models.Post, error)
	GetAllPosts() ([]models.Post, error)
	GetPostById(id uuid.UUID, loggedUserId *uuid.UUID) (models.Post, int8, error)
	UpdatePost(post *models.Post) (models.Post, error)
	DeletePost(post *models.Post) (models.Post, error)
	GetPostsByUserId(userId uuid.UUID, loggedUserId *uuid.UUID) ([]models.Post, map[uuid.UUID]int8, error)
	PostVote(postId uuid.UUID, userId uuid.UUID, vote int8) (models.PostVote, error)
	GetTopPostsOfTheWeek(limit int) ([]models.Post, error)
}

type PostService struct {
	Repo    PostRepositoryInterface
	Storage storage.StorageService
}

func (s *PostService) resolvePostMedia(p *models.Post) {
	if p == nil || s.Storage == nil {
		return
	}
	p.Image = s.Storage.ResolveURL(p.Image)
	p.User.AvatarURL = s.Storage.ResolveURL(p.User.AvatarURL)
	p.User.ProfilePicture = s.Storage.ResolveURL(p.User.ProfilePicture)
	p.User.BannerURL = s.Storage.ResolveURL(p.User.BannerURL)
}

func (s *PostService) UploadPostImage(ctx context.Context, userID uuid.UUID, fileHeader *multipart.FileHeader) (string, string, error) {
	if s.Storage == nil {
		return "", "", apperrors.NewInternalError("almacenamiento multimedia no configurado", nil)
	}

	objectKey, err := s.Storage.SaveImage(ctx, fileHeader, storage.SaveOptions{
		Kind:     storage.PostKind,
		OwnerID:  userID.String(),
		MaxBytes: maxPostImageBytes,
	})
	if err != nil {
		return "", "", apperrors.NewBadRequest(err.Error())
	}

	return objectKey, s.Storage.ResolveURL(objectKey), nil
}

func (s *PostService) CreatePost(post models.Post) (models.Post, error) {
	if post.Image != "" && !storage.IsManagedObjectKey(post.Image) {
		return models.Post{}, apperrors.NewBadRequest("referencia de imagen de publicación inválida")
	}
	created, err := s.Repo.CreatePost(post)
	if err != nil {
		return created, err
	}
	s.resolvePostMedia(&created)
	return created, nil
}

func (s *PostService) GetAllPosts() ([]models.Post, error) {
	posts, err := s.Repo.GetAllPosts()
	if err != nil {
		return nil, err
	}
	for i := range posts {
		s.resolvePostMedia(&posts[i])
	}
	return posts, nil
}

func (s *PostService) GetPostById(id uuid.UUID, loggedUserId *uuid.UUID) (models.Post, int8, error) {
	post, vote, err := s.Repo.GetPostById(id, loggedUserId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Post{}, 0, apperrors.ErrPostNotFound
		}
		return models.Post{}, 0, apperrors.NewInternalError("error al obtener post", err)
	}
	s.resolvePostMedia(&post)
	return post, vote, nil
}

func (s *PostService) UpdatePost(post models.Post) (models.Post, error) {
	updatedPost, _, err := s.Repo.GetPostById(post.ID, nil)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Post{}, apperrors.ErrPostNotFound
		}
		return models.Post{}, apperrors.NewInternalError("error al buscar post", err)
	}

	previousImageKey := updatedPost.Image
	if post.Image != "" && !storage.IsManagedObjectKey(post.Image) {
		return models.Post{}, apperrors.NewBadRequest("referencia de imagen de publicación inválida")
	}

	updatedPost.Title = post.Title
	updatedPost.Content = post.Content
	updatedPost.Image = post.Image
	updatedPost.Type = post.Type
	updatedPost.Tags = post.Tags
	updatedPost.Status = enums.PostStatus(string(post.Status))
	if post.MediaURL != "" {
		updatedPost.MediaURL = post.MediaURL
	}

	saved, err := s.Repo.UpdatePost(&updatedPost)
	if err != nil {
		return models.Post{}, err
	}

	// Eliminar objeto anterior en S3 si fue reemplazado o removido
	if s.Storage != nil && previousImageKey != "" && previousImageKey != post.Image && storage.IsManagedObjectKey(previousImageKey) {
		if delErr := s.Storage.DeleteObject(context.Background(), previousImageKey); delErr != nil {
			log.Printf("warning: failed to delete replaced post image %s: %v", previousImageKey, delErr)
		}
	}

	s.resolvePostMedia(&saved)
	return saved, nil
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

	previousImageKey := post.Image
	post.Status = enums.DELETED
	deleted, err := s.Repo.DeletePost(&post)
	if err != nil {
		return models.Post{}, err
	}

	if s.Storage != nil && previousImageKey != "" && storage.IsManagedObjectKey(previousImageKey) {
		if delErr := s.Storage.DeleteObject(context.Background(), previousImageKey); delErr != nil {
			log.Printf("warning: failed to delete post image on post delete %s: %v", previousImageKey, delErr)
		}
	}

	return deleted, nil
}

// Devuelve los posts de un usuario y el voto del usuario logueado (si se pasa loggedUserId)
func (s *PostService) GetPostsByUserId(userId uuid.UUID, loggedUserId *uuid.UUID) ([]models.Post, map[uuid.UUID]int8, error) {
	posts, votes, err := s.Repo.GetPostsByUserId(userId, loggedUserId)
	if err != nil {
		return nil, nil, err
	}
	for i := range posts {
		s.resolvePostMedia(&posts[i])
	}
	return posts, votes, nil
}

func (s *PostService) PostVote(postId uuid.UUID, userId uuid.UUID, vote int8) (models.PostVote, error) {
	postvote, err := s.Repo.PostVote(postId, userId, vote)
	if err != nil {
		return models.PostVote{}, err
	}
	return postvote, nil
}

// GetCartelera devuelve los posts más populares de la semana
func (s *PostService) GetCartelera(limit int) ([]models.Post, error) {
	if limit <= 0 {
		limit = 10
	}
	posts, err := s.Repo.GetTopPostsOfTheWeek(limit)
	if err != nil {
		return nil, err
	}
	for i := range posts {
		s.resolvePostMedia(&posts[i])
	}
	return posts, nil
}
