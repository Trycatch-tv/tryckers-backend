package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/config"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func setupMiddlewareRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"userId": c.GetString("userId"),
			"role":   c.GetString("role"),
		})
	})
	return r
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	r := setupMiddlewareRouter()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	r := setupMiddlewareRouter()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer ")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestAuthMiddleware_RejectsRefreshToken(t *testing.T) {
	refreshToken, err := utils.CreateRefreshToken("550e8400-e29b-41d4-a716-446655440000", enums.Member)
	if err != nil {
		t.Fatalf("failed to create refresh token: %v", err)
	}

	r := setupMiddlewareRouter()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+refreshToken)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for refresh token, got %d", rec.Code)
	}
}

func TestAuthMiddleware_MalformedClaimsDoesNotPanic(t *testing.T) {
	secret := []byte(config.Load().JWT_SECRET)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  time.Now().Add(time.Hour).Unix(),
		"type": "access",
	})
	signed, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	r := setupMiddlewareRouter()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+signed)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing claims, got %d", rec.Code)
	}
}

func TestAuthMiddleware_ValidAccessToken(t *testing.T) {
	accessToken, err := utils.CreateToken("550e8400-e29b-41d4-a716-446655440000", enums.Member)
	if err != nil {
		t.Fatalf("failed to create access token: %v", err)
	}

	r := setupMiddlewareRouter()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
