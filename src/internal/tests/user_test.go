package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/dtos"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {

	var user = dtos.CreateUserDTO{
		Name:     "test",
		Country:  "colombia",
		Email:    "test@gmail.com",
		Password: "passwordTesting123456789",
	}

	var router = SetupTestRouter()

	t.Run("TestRegisterUser", func(t *testing.T) {
		w := httptest.NewRecorder()

		body := EncodeJSON(user)
		req, _ := http.NewRequest("POST", *GetBaseRoute()+"/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		//req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		response, err := DecodeJSON[models.User](w)

		assert.NoError(t, err)
		assert.Equal(t, user.Name, response.Name)
		assert.Equal(t, user.Email, response.Email)
	})

	t.Run("TestLoginUser", func(t *testing.T) {
		w := httptest.NewRecorder()

		var loginUser = dtos.LoginUser{
			Email:    user.Email,
			Password: user.Password,
		}

		body := EncodeJSON(loginUser)
		req, _ := http.NewRequest("POST", *GetBaseRoute()+"/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		//todo: cuando se haga el merge de la pr de jwt toca actulizar el type LoginResponse
		type LoginResponse struct {
			User models.User `json:"user"`
			//token string
		}

		response, err := DecodeJSON[LoginResponse](w)

		assert.NoError(t, err)
		assert.Equal(t, user.Name, response.User.Name)
		assert.Equal(t, user.Email, response.User.Email)
	})

	t.Run("TestGetAll", func(t *testing.T) {
		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", *GetBaseRoute()+"/users", nil)

		router.ServeHTTP(w, req)

		response, err := DecodeJSON[[]models.User](w)

		assert.NoError(t, err)
		assert.IsType(t, []models.User{}, response)

		found := false
		for _, u := range response {
			if u.Name == user.Name && u.Email == user.Email {
				found = true
				break
			}
		}
		assert.True(t, found, "Not found user")

	})

}
