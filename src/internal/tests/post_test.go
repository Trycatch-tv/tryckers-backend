package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dtos "github.com/Trycatch-tv/tryckers-backend/src/internal/dtos/post"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestPosts(t *testing.T) {
	var router = SetupTestRouter()

	var user = HelperCreateUser(t, router)

	var postDto = dtos.CreatePostDto{
		Title:   "post test" + fmt.Sprint(time.Now().UnixNano()),
		Content: "este es el post test en tryckers",
		Status:  "published",
		UserId:  user.ID,
	}

	t.Run("TestCreatePost", func(t *testing.T) {

		w := httptest.NewRecorder()

		body := EncodeJSON(postDto)
		req, _ := http.NewRequest("POST", *GetBaseRoute()+"/posts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		response, err := DecodeJSON[models.Post](w)

		assert.NoError(t, err)
		assert.Equal(t, postDto.Title, response.Title)
		assert.Equal(t, postDto.UserId, response.User.ID)
	})

	t.Run("TestGetAllPosts", func(t *testing.T) {

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", *GetBaseRoute()+"/posts", nil)

		router.ServeHTTP(w, req)

		response, err := DecodeJSON[[]models.Post](w)

		assert.NoError(t, err)
		assert.IsType(t, []models.Post{}, response)

		found := false
		for _, u := range response {
			if u.Title == postDto.Title {
				found = true
				break
			}
		}
		assert.True(t, found, "Not found post")
	})

	t.Run("TestGetPostById", func(t *testing.T) {

		var createdPost = HelperCreatePost(t, router)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", *GetBaseRoute()+"/posts/"+createdPost.ID.String(), nil)

		router.ServeHTTP(w, req)

		response, err := DecodeJSON[models.Post](w)
		fmt.Print(response)
		assert.NoError(t, err)
		assert.Equal(t, createdPost.ID, response.ID)
	})

	t.Run("TestUpdatePost", func(t *testing.T) {

		var createdPost = HelperCreatePost(t, router)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("PUT", *GetBaseRoute()+"/posts/"+createdPost.ID.String(), nil)

		router.ServeHTTP(w, req)

		response, err := DecodeJSON[models.Post](w)
		fmt.Print(response)
		assert.NoError(t, err)
		assert.Equal(t, createdPost.ID, response.ID)
	})
}
