package middlewares

import (
	"net/http"
	"strings"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/gin-gonic/gin"
)

func RoleMiddleware(role ...enums.UserRole) gin.HandlerFunc {
	println("RoleMiddleware initialized with role:", role)
	return func(ctx *gin.Context) {
		roleToken, exists := ctx.Get("role")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "role not found",
			})
			return
		}

		roleTokenStr := roleToken.(string)
		normalizedRoleToken := strings.ToLower(strings.TrimSpace(roleTokenStr))

		authorized := false
		for _, r := range role {
			normalizedRole := strings.ToLower(strings.TrimSpace(string(r)))
			if normalizedRoleToken == normalizedRole {
				authorized = true
				break
			}
		}

		if !authorized {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User not unauthorized ",
			})
			return
		}

		ctx.Next()

	}

}
