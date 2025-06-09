package services

import (
	"errors"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/dtos"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repository.UserRepository
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.Repo.GetAll()
}

func (s *UserService) CreateUser(user *dtos.CreateUserDTO) (models.User, error) {
	return s.Repo.CreateUser(user)
}

func (s *UserService) Login(user *dtos.LoginUser) (string, error) {

	userData, err := s.Repo.FindByEmail(user)

	if err != nil {
		return "", errors.New("email o contraseña incorrecta")
	}

	IsAuthenticated := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(user.Password))
	if IsAuthenticated != nil {

		return "", errors.New("email o contraseña incorrecta")
	}

	return "token", nil //todo: retornar los datos del usuario y el token
}
