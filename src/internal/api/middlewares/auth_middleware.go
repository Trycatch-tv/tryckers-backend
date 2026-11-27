package middlewares

import (
	"net/http"
	"strings"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" || strings.TrimSpace(parts[1]) == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header must be in format 'Bearer <token>'",
			})
			return
		}

		token := parts[1]
		claims, err := utils.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token invalid",
			})
			return
		}

		if tokenType, ok := claims["type"].(string); !ok || tokenType != "access" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token invalid",
			})
			return
		}

		userId, okUser := claims["sub"].(string)
		role, okRole := claims["role"].(string)
		if !okUser || strings.TrimSpace(userId) == "" || !okRole || strings.TrimSpace(role) == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token invalid",
			})
			return
		}

		ctx.Set("userId", userId)
		ctx.Set("role", role)

		ctx.Next()
	}
}
