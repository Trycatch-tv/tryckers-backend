package tests

import (
	"encoding/json"
	"log"
	"net/http/httptest"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/api/routes"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/utils"
	"github.com/gin-gonic/gin"
)

func SetupTestRouter() *gin.Engine {

	r := gin.Default()
	routes.SetupV1(r, Testdb)

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

func GenerateTokenAdmin() string {
	var id = "id"
	var role = enums.Admin
	token, err := utils.CreateToken(id, role)

	if err != nil {
		log.Fatal("❌ error generating admin token", err)
		return ""
	}

	return token
}

func GenerateTokenMember() string {
	var id = "id"
	var role = enums.Member
	token, err := utils.CreateToken(id, role)

	if err != nil {
		log.Fatal("❌ error generating member token", err)
		return ""
	}

	return token
}
