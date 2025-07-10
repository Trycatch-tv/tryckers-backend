package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
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
		}

		fmt.Println(claims)
		userId := claims["sub"].(string)
		role := claims["role"].(string)

		ctx.Set("userId", userId)
		ctx.Set("role", role)

		ctx.Next()

	}

}
