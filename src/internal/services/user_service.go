package services

import (
	"strings"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/dtos"
	enums "github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	apperrors "github.com/Trycatch-tv/tryckers-backend/src/internal/errors"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/repository"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/utils"
)

type UserService struct {
	Repo *repository.UserRepository
}

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
