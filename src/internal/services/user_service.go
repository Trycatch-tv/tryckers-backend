package services

import (
	"context"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/dtos"
	enums "github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	apperrors "github.com/Trycatch-tv/tryckers-backend/src/internal/errors"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/services/storage"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/utils"
	"github.com/google/uuid"
)

type UserRepositoryInterface interface {
	GetAll() ([]models.User, error)
	CreateUser(user *models.User) (models.User, error)
	FindByEmail(email string) (models.User, error)
	FindByUsername(username string) (models.User, error)
	FindByID(id uuid.UUID) (models.User, error)
	UpdateAvatarURL(id uuid.UUID, avatarURL string) (models.User, error)
	UpdateBannerURL(id uuid.UUID, bannerURL string) (models.User, error)
	UpdateProfile(id uuid.UUID, updates map[string]interface{}) (models.User, error)
}

type UserService struct {
	Repo    UserRepositoryInterface
	Storage storage.StorageService
}

func (s *UserService) resolveUserMedia(u *models.User) {
	if u == nil || s.Storage == nil {
		return
	}
	u.AvatarURL = s.Storage.ResolveURL(u.AvatarURL)
	u.ProfilePicture = s.Storage.ResolveURL(u.ProfilePicture)
	u.BannerURL = s.Storage.ResolveURL(u.BannerURL)
}

const (
	maxAvatarBytes = 2 * 1024 * 1024
	maxBannerBytes = 5 * 1024 * 1024
)

func (s *UserService) GetAllUsers() ([]models.User, error) {
	users, err := s.Repo.GetAll()
	if err != nil {
		return nil, err
	}
	for i := range users {
		s.resolveUserMedia(&users[i])
	}
	return users, nil
}

func (s *UserService) CreateUser(user *dtos.CreateUserDTO) (models.User, error) {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return models.User{}, err
	}
	// Username canónico: sin espacios al borde y en minúsculas. Garantiza que la unicidad
	// no admita duplicados que difieran solo en mayúsculas y que la búsqueda case-insensitive
	// de perfil siempre coincida.
	normalizedUsername := strings.ToLower(strings.TrimSpace(user.Username))
	newUser := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Username: normalizedUsername,
		Password: string(hashedPassword),
		Role:     enums.Member,
		Points:   0,
		Country:  enums.Country(user.Country),
	}

	createdUser, err := s.Repo.CreateUser(&newUser)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "uni_users_username") || strings.Contains(errMsg, "idx_users_username_lower_unique") {
			return models.User{}, apperrors.ErrDuplicateUsername
		}
		if strings.Contains(errMsg, "uni_users_email") {
			return models.User{}, apperrors.ErrDuplicateEmail
		}
		return models.User{}, apperrors.NewInternalError("error al crear usuario", err)
	}
	return createdUser, nil
}

func (s *UserService) Login(user *dtos.LoginUser) (dtos.LoginResponse, error) {
	userData, err := s.Repo.FindByEmail(user.Email)
	if err != nil {
		return dtos.LoginResponse{}, apperrors.ErrInvalidCredentials
	}

	IsAuthenticated := utils.ComparePassword(userData.Password, user.Password)
	if !IsAuthenticated {
		return dtos.LoginResponse{}, apperrors.ErrInvalidCredentials
	}

	token, err := utils.CreateToken(userData.ID.String(), userData.Role)
	if err != nil {
		return dtos.LoginResponse{}, apperrors.NewInternalError("error al generar token", err)
	}

	refreshToken, err := utils.CreateRefreshToken(userData.ID.String(), userData.Role)
	if err != nil {
		return dtos.LoginResponse{}, apperrors.NewInternalError("error al generar refresh token", err)
	}

	s.resolveUserMedia(&userData)

	return dtos.LoginResponse{
		UserData:     userData,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) Perfil(username string) (models.User, error) {
	userPerfil, err := s.Repo.FindByUsername(username)
	if err != nil {
		return models.User{}, apperrors.ErrUserNotFound
	}
	s.resolveUserMedia(&userPerfil)
	return userPerfil, nil
}

func (s *UserService) UploadAvatar(userID uuid.UUID, fileHeader *multipart.FileHeader) (models.User, error) {
	ctx := context.Background()
	user, err := s.Repo.FindByID(userID)
	if err != nil {
		return models.User{}, apperrors.ErrUserNotFound
	}

	objectKey, err := s.Storage.SaveImage(ctx, fileHeader, storage.SaveOptions{
		Kind:     storage.AvatarKind,
		OwnerID:  userID.String(),
		MaxBytes: maxAvatarBytes,
	})
	if err != nil {
		return models.User{}, apperrors.NewBadRequest(err.Error())
	}

	updatedUser, err := s.Repo.UpdateAvatarURL(userID, objectKey)
	if err != nil {
		if delErr := s.Storage.DeleteObject(ctx, objectKey); delErr != nil {
			log.Printf("warning: failed to rollback avatar object %s: %v", objectKey, delErr)
		}
		return models.User{}, apperrors.NewInternalError("error al actualizar avatar", err)
	}

	if user.AvatarURL != "" && user.AvatarURL != objectKey {
		if delErr := s.Storage.DeleteObject(ctx, user.AvatarURL); delErr != nil {
			log.Printf("warning: failed to delete previous avatar %s: %v", user.AvatarURL, delErr)
		}
	}
	if user.ProfilePicture != "" && user.ProfilePicture != user.AvatarURL && user.ProfilePicture != objectKey {
		if delErr := s.Storage.DeleteObject(ctx, user.ProfilePicture); delErr != nil {
			log.Printf("warning: failed to delete previous profile picture %s: %v", user.ProfilePicture, delErr)
		}
	}

	s.resolveUserMedia(&updatedUser)
	return updatedUser, nil
}

func (s *UserService) UploadBanner(userID uuid.UUID, fileHeader *multipart.FileHeader) (models.User, error) {
	ctx := context.Background()
	user, err := s.Repo.FindByID(userID)
	if err != nil {
		return models.User{}, apperrors.ErrUserNotFound
	}

	objectKey, err := s.Storage.SaveImage(ctx, fileHeader, storage.SaveOptions{
		Kind:     storage.BannerKind,
		OwnerID:  userID.String(),
		MaxBytes: maxBannerBytes,
	})
	if err != nil {
		return models.User{}, apperrors.NewBadRequest(err.Error())
	}

	updatedUser, err := s.Repo.UpdateBannerURL(userID, objectKey)
	if err != nil {
		if delErr := s.Storage.DeleteObject(ctx, objectKey); delErr != nil {
			log.Printf("warning: failed to rollback banner object %s: %v", objectKey, delErr)
		}
		return models.User{}, apperrors.NewInternalError("error al actualizar banner", err)
	}

	if user.BannerURL != "" && user.BannerURL != objectKey {
		if delErr := s.Storage.DeleteObject(ctx, user.BannerURL); delErr != nil {
			log.Printf("warning: failed to delete previous banner %s: %v", user.BannerURL, delErr)
		}
	}

	s.resolveUserMedia(&updatedUser)
	return updatedUser, nil
}

func (s *UserService) RemoveAvatar(userID uuid.UUID) (models.User, error) {
	ctx := context.Background()
	user, err := s.Repo.FindByID(userID)
	if err != nil {
		return models.User{}, apperrors.ErrUserNotFound
	}

	updatedUser, err := s.Repo.UpdateAvatarURL(userID, "")
	if err != nil {
		return models.User{}, apperrors.NewInternalError("error al remover avatar", err)
	}

	if user.AvatarURL != "" {
		_ = s.Storage.DeleteObject(ctx, user.AvatarURL)
	}
	if user.ProfilePicture != "" && user.ProfilePicture != user.AvatarURL {
		_ = s.Storage.DeleteObject(ctx, user.ProfilePicture)
	}

	s.resolveUserMedia(&updatedUser)
	return updatedUser, nil
}

func (s *UserService) RemoveBanner(userID uuid.UUID) (models.User, error) {
	ctx := context.Background()
	user, err := s.Repo.FindByID(userID)
	if err != nil {
		return models.User{}, apperrors.ErrUserNotFound
	}

	updatedUser, err := s.Repo.UpdateBannerURL(userID, "")
	if err != nil {
		return models.User{}, apperrors.NewInternalError("error al remover banner", err)
	}

	if user.BannerURL != "" {
		_ = s.Storage.DeleteObject(ctx, user.BannerURL)
	}

	s.resolveUserMedia(&updatedUser)
	return updatedUser, nil
}

func (s *UserService) UpdateProfile(userID uuid.UUID, req *dtos.UpdateProfileDTO) (models.User, error) {
	user, err := s.Repo.FindByID(userID)
	if err != nil {
		return models.User{}, apperrors.ErrUserNotFound
	}

	updates := make(map[string]interface{})

	if req.BirthDate != nil {
		if req.BirthDate.After(time.Now()) {
			return models.User{}, apperrors.NewBadRequest("la fecha de nacimiento no puede ser futura")
		}
		updates["birth_date"] = req.BirthDate
	}

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" {
			return models.User{}, apperrors.NewBadRequest("el nombre no puede estar vacío")
		}
		updates["name"] = trimmed
	}

	if req.Country != nil {
		normalizedCountry := enums.Country(strings.ToLower(strings.TrimSpace(string(*req.Country))))
		if !enums.IsValidCountry(string(normalizedCountry)) {
			return models.User{}, apperrors.NewBadRequest("país inválido")
		}
		updates["country"] = normalizedCountry
	}

	if req.Headline != nil {
		updates["headline"] = strings.TrimSpace(*req.Headline)
	}
	if req.Bio != nil {
		updates["bio"] = strings.TrimSpace(*req.Bio)
	}
	if req.GithubURL != nil {
		updates["github_url"] = strings.TrimSpace(*req.GithubURL)
	}
	if req.LinkedinURL != nil {
		updates["linkedin_url"] = strings.TrimSpace(*req.LinkedinURL)
	}
	if req.PitchVideo != nil {
		updates["pitch_video"] = strings.TrimSpace(*req.PitchVideo)
	}
	if req.Seniority != nil {
		updates["seniority"] = strings.TrimSpace(*req.Seniority)
	}
	if req.EnglishLevel != nil {
		updates["english_level"] = strings.TrimSpace(*req.EnglishLevel)
	}
	if req.EFSetScore != nil {
		updates["ef_set_score"] = strings.TrimSpace(*req.EFSetScore)
	}
	if req.Availability != nil {
		updates["availability"] = strings.TrimSpace(*req.Availability)
	}
	if req.Interests != nil {
		updates["interests"] = strings.TrimSpace(*req.Interests)
	}

	if len(updates) == 0 {
		s.resolveUserMedia(&user)
		return user, nil
	}

	updatedUser, err := s.Repo.UpdateProfile(userID, updates)
	if err != nil {
		return models.User{}, apperrors.NewInternalError("error al actualizar perfil", err)
	}

	s.resolveUserMedia(&updatedUser)
	return updatedUser, nil
}

