package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dtos "github.com/Trycatch-tv/tryckers-backend/src/internal/dtos/post"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type stubPostRepo struct {
	carteleraPosts []models.Post
	carteleraErr   error
}

func (s *stubPostRepo) CreatePost(post models.Post) (models.Post, error) { return post, nil }
func (s *stubPostRepo) GetAllPosts() ([]models.Post, error)              { return nil, nil }
func (s *stubPostRepo) GetPostById(id uuid.UUID, loggedUserId *uuid.UUID) (models.Post, int8, error) {
	return models.Post{}, 0, nil
}
func (s *stubPostRepo) UpdatePost(post *models.Post) (models.Post, error) { return *post, nil }
func (s *stubPostRepo) DeletePost(post *models.Post) (models.Post, error) { return *post, nil }
func (s *stubPostRepo) GetPostsByUserId(userId uuid.UUID, loggedUserId *uuid.UUID) ([]models.Post, map[uuid.UUID]int8, error) {
	return nil, nil, nil
}
func (s *stubPostRepo) PostVote(postId uuid.UUID, userId uuid.UUID, vote int8) (models.PostVote, error) {
	return models.PostVote{}, nil
}
func (s *stubPostRepo) GetTopPostsOfTheWeek(limit int) ([]models.Post, error) {
	return s.carteleraPosts, s.carteleraErr
}

func TestPostHandler_Cartelera_WithPosts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now()
	authorID := uuid.New()
	postID := uuid.New()

	repo := &stubPostRepo{
		carteleraPosts: []models.Post{
			{
				ID:         postID,
				Title:      "Post popular",
				Content:    "Contenido de prueba",
				Type:       enums.RegularPost,
				Status:     enums.PUBLISHED,
				CreatedAt:  &now,
				UpdatedAt:  &now,
				UserID:     authorID,
				VotesCount: 12,
				User: models.User{
					ID:       authorID,
					Username: "devuser",
				},
			},
		},
	}

	handler := &PostHandler{Service: &services.PostService{Repo: repo}}
	r := gin.New()
	r.GET("/api/v1/cartelera", handler.Cartelera)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cartelera", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var got []dtos.ResponsePostDto
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(got) != 1 || got[0].Title != "Post popular" || got[0].VotesCount != 12 {
		t.Fatalf("unexpected cartelera payload: %+v", got)
	}
	if got[0].Author == nil || got[0].Author.Username != "devuser" {
		t.Fatalf("expected author devuser, got %+v", got[0].Author)
	}
}

func TestPostHandler_Cartelera_EmptyReturnsEmptyArray(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &stubPostRepo{carteleraPosts: []models.Post{}}
	handler := &PostHandler{Service: &services.PostService{Repo: repo}}

	r := gin.New()
	r.GET("/api/v1/cartelera", handler.Cartelera)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cartelera", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "[]" {
		t.Fatalf("expected empty JSON array [], got %s", rec.Body.String())
	}
}

func TestPostHandler_Cartelera_DatabaseErrorReturns500WithoutLeakingDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &stubPostRepo{
		carteleraErr: errors.New("pq: relation post_votes does not exist"),
	}
	handler := &PostHandler{Service: &services.PostService{Repo: repo}}

	r := gin.New()
	r.GET("/api/v1/cartelera", handler.Cartelera)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cartelera", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
	var errResp APIError
	if err := json.Unmarshal(rec.Body.Bytes(), &errResp); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if errResp.Error != "error interno del servidor" {
		t.Fatalf("expected generic error message, got %q", errResp.Error)
	}
}
