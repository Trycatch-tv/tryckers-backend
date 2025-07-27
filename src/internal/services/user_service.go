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

	newUser := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hashedPassword),
		Role:     enums.Member,
		Points:   0,
		Country:  enums.Country(user.Country),
	}

	return s.Repo.CreateUser(&newUser)
}

func (s *UserService) Login(user *dtos.LoginUser) (dtos.LoginResponse, error) {

	userData, err := s.Repo.FindByEmail(user.Email)

	if err != nil {
		return dtos.LoginResponse{}, errors.New("incorrect credentials")
	}

	IsAuthenticated := utils.ComparePassword(userData.Password, user.Password)
	if !IsAuthenticated {

		return dtos.LoginResponse{}, errors.New("incorrect credentials")
	}

	token, err := utils.CreateToken(userData.ID.String(), userData.Role)

	if err != nil {
		return dtos.LoginResponse{}, errors.New("internal Server error")
	}

	return dtos.LoginResponse{
		UserData: userData,
		Token:    token}, nil
}

func (s *UserService) Perfil(email string) (models.User, error) {

	userPerfil, err := s.Repo.FindByEmail(email)
	if err != nil {
		return models.User{}, errors.New("user not found")
	}

	return userPerfil, nil
}

func (s *UserService) IsvalidEmail(email string) bool {
	isEmailRegistered := s.Repo.IsEmailRegistered(email)

	return !isEmailRegistered
}
