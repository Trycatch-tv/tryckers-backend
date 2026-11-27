package repository_test

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/config"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/dtos"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	apperrors "github.com/Trycatch-tv/tryckers-backend/src/internal/errors"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/migrations"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/repository"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/services"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openIntegrationDB(t *testing.T) *gorm.DB {
	t.Helper()
	cfg := config.Load()
	dsn := cfg.DBUrl
	if dsn == "" {
		host := cfg.POSTGRES_HOST
		port := cfg.POSTGRES_PORT
		user := cfg.POSTGRES_USER
		pass := cfg.POSTGRES_PASSWORD
		dbName := cfg.POSTGRES_DB
		if user == "" || dbName == "" {
			if envDSN := os.Getenv("TEST_DATABASE_URL"); envDSN != "" {
				dsn = envDSN
			} else {
				t.Skip("PostgreSQL integration DB not configured in environment")
			}
		} else {
			dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbName)
		}
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("PostgreSQL not reachable for integration test: %v", err)
	}
	return db
}

func TestPostgres_UsernameLookupNormalizationAndCaseInsensitiveUniqueness(t *testing.T) {
	db := openIntegrationDB(t)
	if err := migrations.EnsureUsernameCIUniqueIndex(db); err != nil {
		t.Fatalf("EnsureUsernameCIUniqueIndex failed: %v", err)
	}

	repo := &repository.UserRepository{DB: db}
	svc := &services.UserService{Repo: repo}

	suffix := uuid.NewString()[:8]
	mixedUsername := "Deiby7h_" + suffix
	lowerUsername := "deiby7h_" + suffix
	upperUsername := "DEIBY7H_" + suffix
	spacedHistorical := "Charles_" + suffix + " "
	newRawUsername := "  NuevoUsuario_" + suffix + "  "
	newExpectedUsername := "nuevousuario_" + suffix

	var createdIDs []uuid.UUID
	t.Cleanup(func() {
		if len(createdIDs) > 0 {
			_ = db.Where("id IN ?", createdIDs).Delete(&models.User{}).Error
		}
	})

	// 1. Insert historical mixed-case user directly (simulating pre-fix record)
	histUser := models.User{
		Name:     "Historical Mixed Case",
		Username: mixedUsername,
		Email:    fmt.Sprintf("hist_%s@test.local", suffix),
		Password: "hashed-password",
		Role:     enums.Member,
		Country:  enums.Colombia,
	}
	if err := db.Create(&histUser).Error; err != nil {
		t.Fatalf("failed to seed historical user: %v", err)
	}
	createdIDs = append(createdIDs, histUser.ID)

	// 2. Insert historical user with trailing space
	spacedUser := models.User{
		Name:     "Historical Spaced",
		Username: spacedHistorical,
		Email:    fmt.Sprintf("spaced_%s@test.local", suffix),
		Password: "hashed-password",
		Role:     enums.Member,
		Country:  enums.Colombia,
	}
	if err := db.Create(&spacedUser).Error; err != nil {
		t.Fatalf("failed to seed spaced user: %v", err)
	}
	createdIDs = append(createdIDs, spacedUser.ID)

	// 3. Lookup works for mixed, lower, and upper case
	for _, query := range []string{mixedUsername, lowerUsername, upperUsername} {
		found, err := repo.FindByUsername(query)
		if err != nil {
			t.Fatalf("expected FindByUsername(%q) to succeed, got error: %v", query, err)
		}
		if found.ID != histUser.ID {
			t.Fatalf("expected user ID %s, got %s", histUser.ID, found.ID)
		}
	}

	// 4. Lookup works for historical trailing-space username
	foundSpaced, err := repo.FindByUsername("charles_" + suffix)
	if err != nil || foundSpaced.ID != spacedUser.ID {
		t.Fatalf("expected spaced historical user to be found, got err=%v", err)
	}

	// 5. Registering lowercase duplicate of existing historical mixed-case user MUST be rejected (AC-04)
	_, dupErr := svc.CreateUser(&dtos.CreateUserDTO{
		Name:     "Duplicate Attempt",
		Username: lowerUsername,
		Country:  enums.Colombia,
		Email:    fmt.Sprintf("dup_%s@test.local", suffix),
		Password: "password123",
	})
	if !errors.Is(dupErr, apperrors.ErrDuplicateUsername) {
		t.Fatalf("expected ErrDuplicateUsername for case-insensitive duplicate, got: %v", dupErr)
	}

	// 6. Registering new user with spaces & uppercase normalizes to lowercase trimmed
	created, err := svc.CreateUser(&dtos.CreateUserDTO{
		Name:     "New User",
		Username: newRawUsername,
		Country:  enums.Colombia,
		Email:    fmt.Sprintf("new_%s@test.local", suffix),
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("expected CreateUser to succeed, got: %v", err)
	}
	createdIDs = append(createdIDs, created.ID)
	if created.Username != newExpectedUsername {
		t.Fatalf("expected normalized username %q, got %q", newExpectedUsername, created.Username)
	}

	// 7. Non-existent user returns ErrRecordNotFound (maps to 404)
	_, notFoundErr := repo.FindByUsername("no_existe_" + suffix)
	if !errors.Is(notFoundErr, gorm.ErrRecordNotFound) {
		t.Fatalf("expected gorm.ErrRecordNotFound for nonexistent user, got: %v", notFoundErr)
	}
}
