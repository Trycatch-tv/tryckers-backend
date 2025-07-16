package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/utils"
)

func HelperCreateUser(t *testing.T) models.User {
	t.Helper()

	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}
	user := models.User{
		Name:     "Test User",
		Email:    fmt.Sprintf("testuser_%d@example.com", time.Now().UnixNano()),
		Password: hashedPassword,
		Country:  "colombia",
	}

	if err := Testdb.Create(&user).Error; err != nil {
		t.Fatalf("Error creating test user in DB: %v", err)
	}

	return user
}

func HelperCreatePost(t *testing.T) models.Post {
	t.Helper()

	user := HelperCreateUser(t)

	post := models.Post{
		Title:   "post test " + fmt.Sprint(time.Now().UnixNano()),
		Content: "este es el post test en tryckers",
		Status:  "published",
		UserID:  user.ID, // Usa el ID del usuario creado
	}

	if err := Testdb.Create(&post).Error; err != nil {
		t.Fatalf("Error creating test post in DB: %v", err)
	}

	return post
}

func HelperCreateComment(t *testing.T) models.Comment {
	t.Helper()

	// Crear post y usuario directamente en la DB
	post := HelperCreatePost(t)
	user := HelperCreateUser(t)

	comment := models.Comment{
		Content: "comment test " + fmt.Sprint(time.Now().UnixNano()),
		UserID:  user.ID,
		PostID:  post.ID,
	}

	if err := Testdb.Create(&comment).Error; err != nil {
		t.Fatalf("Error creating test comment in DB: %v", err)
	}

	return comment
}
