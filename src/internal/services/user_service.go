package services

import (
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
	newPost := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Username: user.Username,
		Password: string(hashedPassword),
		Role:     enums.Member,
		Points:   0,
		Country:  enums.Country(user.Country),
	}

	return s.Repo.CreateUser(&newPost)
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

	return dtos.LoginResponse{
		UserData: userData,
		Token:    token,
	}, nil
}

func (s *UserService) Perfil(username string) (models.User, error) {
	userPerfil, err := s.Repo.FindByUsername(username)
	if err != nil {
		return models.User{}, apperrors.ErrUserNotFound
	}

	return userPerfil, nil
}
