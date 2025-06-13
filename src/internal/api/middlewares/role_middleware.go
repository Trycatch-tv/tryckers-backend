package middlewares

import (
	"net/http"
	"strings"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/gin-gonic/gin"
)

func RoleMiddleware(role enums.UserRole) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		roleToken, exists := ctx.Get("role")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "role not found",
			})
			return
		}

		// Convertir enums.UserRole a string
		roleStr := string(role) // si enums.UserRole es un tipo alias de string
		roleTokenStr := roleToken.(string)

		// Normalizar ambos
		normalizedRoleToken := strings.ToLower(strings.TrimSpace(roleTokenStr))
		normalizedRole := strings.ToLower(strings.TrimSpace(roleStr))

		if normalizedRoleToken != normalizedRole {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User not unauthorized ",
			})
			return
		}

		ctx.Next()

	}

}
