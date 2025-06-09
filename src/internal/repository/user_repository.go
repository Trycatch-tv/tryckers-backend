package repository

import (
	"github.com/Trycatch-tv/tryckers-backend/src/internal/dtos"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.DB.Find(&users).Error
	return users, err
}

func (r *UserRepository) CreateUser(user *dtos.CreateUserDTO) (models.User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	userModel := models.User{
		Name:     user.Name,
		Email:    user.Email,
		Password: string(hashedPassword),
	}

	result := r.DB.Create(&userModel)
	return userModel, result.Error
}

func (r *UserRepository) FindByEmail(user *dtos.LoginUser) (models.User, error) {
	var foundUser models.User

	result := r.DB.Where("email = ?", user.Email).First(&foundUser)
	return foundUser, result.Error
}
