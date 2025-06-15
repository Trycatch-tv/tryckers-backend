package services

import (
	"errors"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/dtos"
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
	user.Password = hashedPassword
	return s.Repo.CreateUser(user)
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

func (s *UserService) Profile(name *string) (models.User, error) {
	userProfile, err := s.Repo.FindByName(name)
	if err != nil {
		return models.User{}, errors.New("user not found")
	}

	return userProfile, nil
}
