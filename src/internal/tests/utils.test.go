package tests

import (
	"encoding/json"
	"log"
	"net/http/httptest"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/api/routes"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/config"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupTestRouter() *gin.Engine {
	cfg := config.LoadTest()
	db := config.InitGormDB(cfg)

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})

	r := gin.Default()
	routes.SetupV1(r, db)

	return r
}

func GetBaseRoute() *string {
	baseRoute := "/api/v1"
	return &baseRoute
}

// EncodeJSON serializa cualquier struct a []byte (JSON)
func EncodeJSON[T any](data T) []byte {
	body, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Could not serialize JSON: %v", err)
	}
	return body
}

// DecodeJSON deserializa un JSON de respuesta a una estructura
func DecodeJSON[T any](w *httptest.ResponseRecorder) (T, error) {
	var target T
	err := json.Unmarshal(w.Body.Bytes(), &target)

	return target, err
}
