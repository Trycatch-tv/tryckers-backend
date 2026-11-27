package services

import (
	"testing"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/dtos"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/google/uuid"
)

// stubUserRepo captura el usuario persistido para verificar la normalización.
type stubUserRepo struct {
	lastCreatedUser *models.User
}

func (s *stubUserRepo) GetAll() ([]models.User, error) { return nil, nil }
func (s *stubUserRepo) CreateUser(user *models.User) (models.User, error) {
	user.ID = uuid.New()
	s.lastCreatedUser = user
	return *user, nil
}
func (s *stubUserRepo) FindByEmail(email string) (models.User, error)       { return models.User{}, nil }
func (s *stubUserRepo) FindByUsername(username string) (models.User, error) { return models.User{}, nil }
func (s *stubUserRepo) FindByID(id uuid.UUID) (models.User, error)          { return models.User{}, nil }
func (s *stubUserRepo) UpdateAvatarURL(id uuid.UUID, avatarURL string) (models.User, error) {
	return models.User{}, nil
}
func (s *stubUserRepo) UpdateBannerURL(id uuid.UUID, bannerURL string) (models.User, error) {
	return models.User{}, nil
}
func (s *stubUserRepo) UpdateProfile(id uuid.UUID, updates map[string]interface{}) (models.User, error) {
	return models.User{}, nil
}

// TestUserService_CreateUser_NormalizesUsername verifica que el registro persista el
// username en forma canónica (sin espacios al borde y en minúsculas), de modo que la
// unicidad no admita duplicados que difieran solo en mayúsculas y que la búsqueda
// case-insensitive siempre coincida.
func TestUserService_CreateUser_NormalizesUsername(t *testing.T) {
	repo := &stubUserRepo{}
	svc := &UserService{Repo: repo}

	_, err := svc.CreateUser(&dtos.CreateUserDTO{
		Name:     "Deibith Alejandro",
		Username: "  Deiby7h322  ",
		Country:  enums.Colombia,
		Email:    "deiby@test.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("CreateUser returned unexpected error: %v", err)
	}
	if repo.lastCreatedUser == nil {
		t.Fatalf("expected repo to receive a user, got nil")
	}
	if got := repo.lastCreatedUser.Username; got != "deiby7h322" {
		t.Fatalf("expected persisted username %q, got %q", "deiby7h322", got)
	}
}

func TestUserService_CreateUser_RejectsWhitespaceInUsername(t *testing.T) {
	repo := &stubUserRepo{}
	svc := &UserService{Repo: repo}

	for _, bad := range []string{"   ", "juan perez", "deiby\tdev"} {
		_, err := svc.CreateUser(&dtos.CreateUserDTO{
			Name:     "Test User",
			Username: bad,
			Country:  enums.Colombia,
			Email:    "test@test.com",
			Password: "password123",
		})
		if err == nil {
			t.Fatalf("expected error for username %q, got nil", bad)
		}
	}
}
