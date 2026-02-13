package repository

import (
	"fmt"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository struct {
	DB *gorm.DB
}

func (r *PostRepository) CreatePost(post models.Post) (models.Post, error) {
	result := r.DB.Create(&post)
	if result.Error != nil {
		return post, fmt.Errorf("error al crear post: %w", result.Error)
	}
	var createdPost models.Post
	err := r.DB.Preload("User").First(&createdPost, post.ID).Error
	if err != nil {
		return createdPost, fmt.Errorf("error al obtener post creado: %w", err)
	}
	return createdPost, nil
}
func (r *PostRepository) GetAllPosts() ([]models.Post, error) {
	var posts []models.Post
	err := r.DB.Model(&models.Post{}).
		Select("posts.*, COALESCE(SUM(CASE WHEN post_votes.vote = 1 THEN 1 ELSE 0 END), 0) AS votes_count").
		Joins("LEFT JOIN post_votes ON post_votes.post_id = posts.id").
		Where("posts.status != ?", enums.DELETED).
		Group("posts.id").
		Preload("User").
		Find(&posts).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener posts: %w", err)
	}
	return posts, nil
}
func (r *PostRepository) GetPostById(id uuid.UUID, loggedUserId *uuid.UUID) (models.Post, int8, error) {
	var post models.Post
	err := r.DB.Model(&models.Post{}).
		Select("posts.*, COALESCE(SUM(CASE WHEN post_votes.vote = 1 THEN 1 ELSE 0 END), 0) AS votes_count").
		Joins("LEFT JOIN post_votes ON post_votes.post_id = posts.id").
		Where("posts.id = ? AND posts.status != ?", id, enums.DELETED).
		Group("posts.id").
		Preload("User").
		First(&post).Error
	if err != nil {
		return post, 0, err
	}

	var userVote int8 = 0
	if loggedUserId != nil {
		var postVote models.PostVote
		voteErr := r.DB.Model(&models.PostVote{}).
			Where("user_id = ? AND post_id = ?", *loggedUserId, id).
			First(&postVote).Error
		if voteErr == nil {
			userVote = postVote.Vote
		}
	}
	return post, userVote, nil
}
func (r *PostRepository) UpdatePost(post *models.Post) (models.Post, error) {
	result := r.DB.Save(post)
	if result.Error != nil {
		return models.Post{}, fmt.Errorf("error al actualizar post: %w", result.Error)
	}
	return *post, nil
}
func (r *PostRepository) DeletePost(post *models.Post) (models.Post, error) {
	result := r.DB.Save(post)
	if result.Error != nil {
		return models.Post{}, fmt.Errorf("error al eliminar post: %w", result.Error)
	}
	return *post, nil
}

// Devuelve los posts de un usuario y, si se pasa loggedUserId, agrega el voto de ese usuario en cada post
func (r *PostRepository) GetPostsByUserId(userId uuid.UUID, loggedUserId *uuid.UUID) ([]models.Post, map[uuid.UUID]int8, error) {
	var posts []models.Post
	err := r.DB.Model(&models.Post{}).
		Select("posts.*, COALESCE(SUM(CASE WHEN post_votes.vote = 1 THEN 1 ELSE 0 END), 0) AS votes_count").
		Joins("LEFT JOIN post_votes ON post_votes.post_id = posts.id").
		Where("posts.user_id = ? AND posts.status != ?", userId, enums.DELETED).
		Group("posts.id").
		Preload("User").
		Find(&posts).Error
	if err != nil {
		return nil, nil, fmt.Errorf("error al obtener posts del usuario: %w", err)
	}
	userVotes := make(map[uuid.UUID]int8)
	if loggedUserId != nil {
		// Obtener el voto del usuario logueado para cada post
		var votes []models.PostVote
		postIDs := make([]uuid.UUID, 0, len(posts))
		for _, p := range posts {
			postIDs = append(postIDs, p.ID)
		}
		if len(postIDs) > 0 {
			r.DB.Model(&models.PostVote{}).
				Where("user_id = ? AND post_id IN ?", *loggedUserId, postIDs).
				Find(&votes)
			for _, v := range votes {
				userVotes[v.PostID] = v.Vote
			}
		}
	}
	return posts, userVotes, nil
}

func (r *PostRepository) PostVote(postId uuid.UUID, userId uuid.UUID, vote int8) (models.PostVote, error) {
	var postVote models.PostVote
	err := r.DB.Preload("User").Preload("Post").Where("post_id = ? AND user_id = ?", postId, userId).First(&postVote).Error
	if err == gorm.ErrRecordNotFound {
		// No existe: crear registro
		postVote = models.PostVote{
			PostID: postId,
			UserID: userId,
			Vote:   int8(vote),
		}
		if err := r.DB.Create(&postVote).Error; err != nil {
			return models.PostVote{}, fmt.Errorf("error al crear voto: %w", err)
		}
		return postVote, nil
	} else if err != nil {
		// Error inesperado
		return models.PostVote{}, fmt.Errorf("error al buscar voto: %w", err)
	}
	// Sí existe: actualizar el voto
	postVote.Vote = int8(vote)
	if err := r.DB.Save(&postVote).Error; err != nil {
		return models.PostVote{}, fmt.Errorf("error al actualizar voto: %w", err)
	}
	return postVote, nil
}
