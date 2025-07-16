package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dtos "github.com/Trycatch-tv/tryckers-backend/src/internal/dtos/comment"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestComments(t *testing.T) {
	var router = SetupTestRouter()

	var user = HelperCreateUser(t)
	var post = HelperCreatePost(t)

	var comment = dtos.CreateCommentDto{
		Content: "comment test" + fmt.Sprint(time.Now().UnixNano()),
		UserId:  user.ID,
		PostId:  post.ID,
	}

	t.Run("TestCreateComment", func(t *testing.T) {

		w := httptest.NewRecorder()

		body := EncodeJSON(comment)
		req, _ := http.NewRequest("POST", *GetBaseRoute()+"/comments", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		response, err := DecodeJSON[models.Comment](w)

		assert.NoError(t, err)
		assert.Equal(t, user.ID, response.UserID)
		assert.Equal(t, post.ID, response.PostID)

	})

	t.Run("TestGetCommentsByPostId", func(t *testing.T) {

		var comment = HelperCreateComment(t)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("GET", *GetBaseRoute()+"/posts/"+comment.PostID.String()+"/comments", nil)

		router.ServeHTTP(w, req)

		response, err := DecodeJSON[[]models.Comment](w)

		assert.NoError(t, err)

		assert.IsType(t, []models.Comment{}, response)

		found := false
		for _, u := range response {
			if u.ID == comment.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "Not found post")

	})

	t.Run("TestUpdateComment", func(t *testing.T) {

		var createdComment = HelperCreateComment(t)

		var updateComment = createdComment

		updateComment.Content = "update content"

		w := httptest.NewRecorder()

		body := EncodeJSON(updateComment)
		req, _ := http.NewRequest("PUT", *GetBaseRoute()+"/comments/"+updateComment.ID.String(), bytes.NewBuffer(body))
		router.ServeHTTP(w, req)

		response, err := DecodeJSON[models.Comment](w)
		assert.NoError(t, err)
		assert.Equal(t, updateComment.ID, response.ID)
		assert.NotEqual(t, createdComment.Content, response.Content)
	})

	t.Run("TestDeleteComments", func(t *testing.T) {

		var createdComment = HelperCreateComment(t)

		w := httptest.NewRecorder()

		req, _ := http.NewRequest("DELETE", *GetBaseRoute()+"/comments/"+createdComment.ID.String(), nil)

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
	})
}
