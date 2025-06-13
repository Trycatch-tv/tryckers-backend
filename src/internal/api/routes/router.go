package routes

import (
	"github.com/Trycatch-tv/tryckers-backend/src/internal/api/handlers"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/api/middlewares"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/repository"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupV1(r *gin.Engine, db *gorm.DB) {
	userRepo := &repository.UserRepository{DB: db}
	userService := &services.UserService{Repo: userRepo}
	userHandler := &handlers.UserHandler{Service: userService}

	api := r.Group("/api/v1")
	{
		api.GET("/users", middlewares.AuthMiddleware(), middlewares.RoleMiddleware(enums.Admin), userHandler.GetAll)
		api.POST("/register", userHandler.CreateUser)
		api.POST("/login", userHandler.Login)
		api.GET("/perfil/:name", middlewares.AuthMiddleware(), userHandler.Perfil)
	}
}
