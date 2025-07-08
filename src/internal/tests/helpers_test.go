package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/dtos"
	dtoPost "github.com/Trycatch-tv/tryckers-backend/src/internal/dtos/post"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/gin-gonic/gin"
)

func HelperCreateUser(t *testing.T, router *gin.Engine) models.User {
	t.Helper() // helper

	user := dtos.CreateUserDTO{
		Name:     "Test User",
		Email:    fmt.Sprintf("testuser_%d@example.com", time.Now().UnixNano()), // email único
		Password: "password123",
		Country:  "colombia",
	}

	body := EncodeJSON(user)

	req, _ := http.NewRequest("POST", *GetBaseRoute()+"/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("The user could not be created, status code: %d, body: %s", w.Code, w.Body.String())
	}

	response, err := DecodeJSON[models.User](w)
	if err != nil {
		t.Fatalf("Error deserializing user response: %v", err)
	}

	return response
}

func HelperCreatePost(t *testing.T, router *gin.Engine) models.Post {
	t.Helper() // helper
	var user = HelperCreateUser(t, router)
	post := dtoPost.CreatePostDto{
		Title:   "post test" + fmt.Sprint(time.Now().UnixNano()),
		Content: "este es el post test en tryckers",
		Status:  "published",
		UserId:  user.ID,
	}

	body := EncodeJSON(post)

	req, _ := http.NewRequest("POST", *GetBaseRoute()+"/posts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("The post could not be created, status code: %d, body: %s", w.Code, w.Body.String())
	}

	response, err := DecodeJSON[models.Post](w)
	if err != nil {
		t.Fatalf("Error deserializing post response: %v", err)
	}

	return response
}
