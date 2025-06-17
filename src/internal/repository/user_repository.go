package repository

import (
	"github.com/Trycatch-tv/tryckers-backend/src/internal/dtos"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
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

func (r *UserRepository) CreateUser(user *models.User) (models.User, error) {
	result := r.DB.Create(&user)
	return *user, result.Error
}

func (r *UserRepository) FindByEmail(user *dtos.LoginUser) (models.User, error) {
	var foundUser models.User

	result := r.DB.Where("email = ?", user.Email).First(&foundUser)
	return foundUser, result.Error
}

func (r *UserRepository) FindByName(name *string) (models.User, error) {
	var foundUser models.User

	result := r.DB.Where("name = ?", name).First(&foundUser)
	return foundUser, result.Error
}
