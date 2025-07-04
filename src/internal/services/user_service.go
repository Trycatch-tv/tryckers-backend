package services

import (
	"errors"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/dtos"
	enums "github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
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
		Password: string(hashedPassword),
		Role:     enums.Member,
		Points:   0,
		Country:  enums.Country(user.Country),
	}

	return s.Repo.CreateUser(&newPost)
}

func (s *UserService) Login(user *dtos.LoginUser) (models.User, error) {

	userData, err := s.Repo.FindByEmail(user)

	if err != nil {
		return models.User{}, errors.New("incorrect credentials")
	}

	IsAuthenticated := utils.ComparePassword(userData.Password, user.Password)
	if !IsAuthenticated {

		return models.User{}, errors.New("incorrect credentials")
	}

	return userData, nil
}

func (s *UserService) Perfil(name *string) (models.User, error) {

	userPerfil, err := s.Repo.FindByName(name)
	if err != nil {
		return models.User{}, errors.New("user not found")
	}

	return userPerfil, nil
}
