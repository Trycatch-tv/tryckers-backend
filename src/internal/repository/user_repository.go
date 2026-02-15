package repository

import (
	"fmt"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.DB.Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("error al obtener usuarios: %w", err)
	}
	return users, nil
}

func (r *UserRepository) CreateUser(user *models.User) (models.User, error) {
	result := r.DB.Create(&user)
	if result.Error != nil {
		return models.User{}, fmt.Errorf("error al crear usuario: %w", result.Error)
	}
	return *user, nil
}

func (r *UserRepository) FindByEmail(email string) (models.User, error) {
	var foundUser models.User
	result := r.DB.Where("email = ?", email).First(&foundUser)
	if result.Error != nil {
		return models.User{}, fmt.Errorf("error al buscar usuario por email: %w", result.Error)
	}
	return foundUser, nil
}

func (r *UserRepository) FindByUsername(username string) (models.User, error) {
	var foundUser models.User
	result := r.DB.Where("username = ?", username).First(&foundUser)
	if result.Error != nil {
		return models.User{}, fmt.Errorf("error al buscar usuario por username: %w", result.Error)
	}
	return foundUser, nil
}
