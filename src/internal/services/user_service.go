package services

import (
	"mime/multipart"
	"strings"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/dtos"
	enums "github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	apperrors "github.com/Trycatch-tv/tryckers-backend/src/internal/errors"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/repository"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/services/storage"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/utils"
	"github.com/google/uuid"
)

type UserService struct {
	Repo    *repository.UserRepository
	Storage *storage.LocalStorage
}

const (
	maxAvatarBytes = 2 * 1024 * 1024
	maxBannerBytes = 5 * 1024 * 1024
)

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.Repo.GetAll()
}

func (s *UserService) CreateUser(user *dtos.CreateUserDTO) (models.User, error) {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return models.User{}, err
	}
	newUser := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
		Password: string(hashedPassword),
		Role:     enums.Member,
		Points:   0,
		Country:  enums.Country(user.Country),
	}

	createdUser, err := s.Repo.CreateUser(&newUser)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "uni_users_username") {
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

	return userPerfil, nil
}

func (s *UserService) UploadAvatar(userID uuid.UUID, fileHeader *multipart.FileHeader) (models.User, error) {
	user, err := s.Repo.FindByID(userID)
	if err != nil {
		return models.User{}, apperrors.ErrUserNotFound
	}

	avatarURL, err := s.Storage.SaveImage(fileHeader, storage.SaveOptions{
		Kind:     storage.AvatarKind,
		UserID:   userID.String(),
		MaxBytes: maxAvatarBytes,
	})
	if err != nil {
		return models.User{}, apperrors.NewBadRequest(err.Error())
	}

	updatedUser, err := s.Repo.UpdateAvatarURL(userID, avatarURL)
	if err != nil {
		s.Storage.DeleteByPublicURL(avatarURL)
		return models.User{}, apperrors.NewInternalError("error al actualizar avatar", err)
	}

	s.Storage.DeleteByPublicURL(user.AvatarURL)
	if user.ProfilePicture != user.AvatarURL {
		s.Storage.DeleteByPublicURL(user.ProfilePicture)
	}

	return updatedUser, nil
}

func (s *UserService) UploadBanner(userID uuid.UUID, fileHeader *multipart.FileHeader) (models.User, error) {
	user, err := s.Repo.FindByID(userID)
	if err != nil {
		return models.User{}, apperrors.ErrUserNotFound
	}

	bannerURL, err := s.Storage.SaveImage(fileHeader, storage.SaveOptions{
		Kind:     storage.BannerKind,
		UserID:   userID.String(),
		MaxBytes: maxBannerBytes,
	})
	if err != nil {
		return models.User{}, apperrors.NewBadRequest(err.Error())
	}

	updatedUser, err := s.Repo.UpdateBannerURL(userID, bannerURL)
	if err != nil {
		s.Storage.DeleteByPublicURL(bannerURL)
		return models.User{}, apperrors.NewInternalError("error al actualizar banner", err)
	}

	s.Storage.DeleteByPublicURL(user.BannerURL)

	return updatedUser, nil
}

func (s *UserService) RemoveAvatar(userID uuid.UUID) (models.User, error) {
	user, err := s.Repo.FindByID(userID)
	if err != nil {
		return models.User{}, apperrors.ErrUserNotFound
	}

	updatedUser, err := s.Repo.UpdateAvatarURL(userID, "")
	if err != nil {
		return models.User{}, apperrors.NewInternalError("error al remover avatar", err)
	}

	s.Storage.DeleteByPublicURL(user.AvatarURL)
	if user.ProfilePicture != user.AvatarURL {
		s.Storage.DeleteByPublicURL(user.ProfilePicture)
	}

	return updatedUser, nil
}

func (s *UserService) RemoveBanner(userID uuid.UUID) (models.User, error) {
	user, err := s.Repo.FindByID(userID)
	if err != nil {
		return models.User{}, apperrors.ErrUserNotFound
	}

	updatedUser, err := s.Repo.UpdateBannerURL(userID, "")
	if err != nil {
		return models.User{}, apperrors.NewInternalError("error al remover banner", err)
	}

	s.Storage.DeleteByPublicURL(user.BannerURL)

	return updatedUser, nil
}
