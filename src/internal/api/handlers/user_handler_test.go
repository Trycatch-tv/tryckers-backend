package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type stubUserRepo struct {
	user            models.User
	lastCreatedUser *models.User
	lastUpdates     map[string]interface{}
}

func (s *stubUserRepo) GetAll() ([]models.User, error) { return []models.User{s.user}, nil }
func (s *stubUserRepo) CreateUser(user *models.User) (models.User, error) {
	user.ID = uuid.New()
	s.lastCreatedUser = user
	return *user, nil
}
func (s *stubUserRepo) FindByEmail(email string) (models.User, error)       { return s.user, nil }
func (s *stubUserRepo) FindByUsername(username string) (models.User, error) { return s.user, nil }
func (s *stubUserRepo) FindByID(id uuid.UUID) (models.User, error)          { return s.user, nil }
func (s *stubUserRepo) UpdateAvatarURL(id uuid.UUID, avatarURL string) (models.User, error) {
	return s.user, nil
}
func (s *stubUserRepo) UpdateBannerURL(id uuid.UUID, bannerURL string) (models.User, error) {
	return s.user, nil
}
func (s *stubUserRepo) UpdateProfile(id uuid.UUID, updates map[string]interface{}) (models.User, error) {
	s.lastUpdates = updates
	updated := s.user
	if bio, ok := updates["bio"].(string); ok {
		updated.Bio = bio
	}
	if country, ok := updates["country"].(enums.Country); ok {
		updated.Country = country
	}
	return updated, nil
}

func TestUserHandler_CreateUser_RejectsMissingCountry(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &stubUserRepo{}
	handler := &UserHandler{Service: &services.UserService{Repo: repo}}
	r := gin.New()
	r.POST("/api/v1/register", handler.CreateUser)

	body := `{"name":"Juan Perez","username":"juanp","email":"juan@test.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 when country is missing, got %d (%s)", rec.Code, rec.Body.String())
	}
}

func TestUserHandler_CreateUser_RejectsInvalidISOCountry(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &stubUserRepo{}
	handler := &UserHandler{Service: &services.UserService{Repo: repo}}
	r := gin.New()
	r.POST("/api/v1/register", handler.CreateUser)

	body := `{"name":"Juan Perez","username":"juanp","country":"CO","email":"juan@test.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for ISO code CO, got %d", rec.Code)
	}
}

func TestUserHandler_CreateUser_AcceptsCanonicalCountry(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &stubUserRepo{}
	handler := &UserHandler{Service: &services.UserService{Repo: repo}}
	r := gin.New()
	r.POST("/api/v1/register", handler.CreateUser)

	body := `{"name":"Juan Perez","username":"juanp","country":"colombia","email":"juan@test.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 for valid country, got %d (%s)", rec.Code, rec.Body.String())
	}
	if repo.lastCreatedUser == nil || repo.lastCreatedUser.Country != enums.Colombia {
		t.Fatalf("expected persisted country colombia, got %+v", repo.lastCreatedUser)
	}
}

func TestUserHandler_UpdateMe_AllowsBioUpdateWithoutCountry(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	repo := &stubUserRepo{
		user: models.User{
			ID:       userID,
			Name:     "Usuario Sin Pais",
			Username: "sinpais",
			Country:  "",
		},
	}
	handler := &UserHandler{Service: &services.UserService{Repo: repo}}
	r := gin.New()
	r.PUT("/api/v1/users/me", func(c *gin.Context) {
		c.Set("userId", userID.String())
		handler.UpdateMe(c)
	})

	body := `{"bio":"Nueva biografia actualizada"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 updating bio without country, got %d (%s)", rec.Code, rec.Body.String())
	}
	if _, hasCountry := repo.lastUpdates["country"]; hasCountry {
		t.Fatalf("expected country not to be in updates map, got %+v", repo.lastUpdates)
	}
}

func TestUserHandler_UpdateMe_RejectsInvalidCountryAndAcceptsCanonical(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userID := uuid.New()
	repo := &stubUserRepo{
		user: models.User{
			ID:       userID,
			Name:     "Usuario Existente",
			Username: "existente",
			Country:  enums.Colombia,
		},
	}
	handler := &UserHandler{Service: &services.UserService{Repo: repo}}
	r := gin.New()
	r.PUT("/api/v1/users/me", func(c *gin.Context) {
		c.Set("userId", userID.String())
		handler.UpdateMe(c)
	})

	// 1. País inválido ("CO") -> 400
	badReq := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", bytes.NewBufferString(`{"country":"CO"}`))
	badReq.Header.Set("Content-Type", "application/json")
	badRec := httptest.NewRecorder()
	r.ServeHTTP(badRec, badReq)
	if badRec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid country CO, got %d", badRec.Code)
	}

	var apiErr APIError
	_ = json.Unmarshal(badRec.Body.Bytes(), &apiErr)
	if apiErr.Error != "país inválido" {
		t.Fatalf("expected specific error 'país inválido', got %q", apiErr.Error)
	}

	// 2. País canónico ("mexico") -> 200
	goodReq := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", bytes.NewBufferString(`{"country":"mexico","bio":"Bio ok"}`))
	goodReq.Header.Set("Content-Type", "application/json")
	goodRec := httptest.NewRecorder()
	r.ServeHTTP(goodRec, goodReq)
	if goodRec.Code != http.StatusOK {
		t.Fatalf("expected 200 for canonical country mexico, got %d (%s)", goodRec.Code, goodRec.Body.String())
	}
	if repo.lastUpdates["country"] != enums.Mexico {
		t.Fatalf("expected country updated to mexico, got %v", repo.lastUpdates["country"])
	}
}
