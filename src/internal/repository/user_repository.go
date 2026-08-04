package repository

import (
	"fmt"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/google/uuid"
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

func (r *UserRepository) FindByID(id uuid.UUID) (models.User, error) {
	var foundUser models.User
	result := r.DB.First(&foundUser, "id = ?", id)
	if result.Error != nil {
		return models.User{}, fmt.Errorf("error al buscar usuario por id: %w", result.Error)
	}
	return foundUser, nil
}

func (r *UserRepository) UpdateAvatarURL(id uuid.UUID, avatarURL string) (models.User, error) {
	updates := map[string]interface{}{
		"avatar_url":      avatarURL,
		"profile_picture": avatarURL,
	}

	if err := r.DB.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return models.User{}, fmt.Errorf("error al actualizar avatar: %w", err)
	}

	return r.FindByID(id)
}

func (r *UserRepository) UpdateBannerURL(id uuid.UUID, bannerURL string) (models.User, error) {
	if err := r.DB.Model(&models.User{}).Where("id = ?", id).Update("banner_url", bannerURL).Error; err != nil {
		return models.User{}, fmt.Errorf("error al actualizar banner: %w", err)
	}

	return r.FindByID(id)
}

func (r *UserRepository) UpdateProfile(id uuid.UUID, updates map[string]interface{}) (models.User, error) {
	if len(updates) == 0 {
		return r.FindByID(id)
	}

	if err := r.DB.Model(&models.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return models.User{}, fmt.Errorf("error al actualizar perfil: %w", err)
	}

	return r.FindByID(id)
}

