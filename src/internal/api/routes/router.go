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
	postRepo := &repository.PostRepository{DB: db}
	postService := &services.PostService{Repo: postRepo}
	postHandler := &handlers.PostHandler{Service: postService}
	commentRepo := &repository.CommentRepository{DB: db}
	commentService := &services.CommentService{Repo: commentRepo}
	commentHandler := &handlers.CommentHandler{Service: commentService, PostService: postService}
	api := r.Group("/api/v1")
	{
		// Rutas públicas (no requieren autenticación)
		api.POST("/register", userHandler.CreateUser)
		api.POST("/login", userHandler.Login)
		api.POST("/refresh-token", userHandler.RefreshToken)

		// Rutas protegidas (requieren autenticación)
		protected := api.Group("/")
		protected.Use(middlewares.AuthMiddleware())
		{
			// Users
			protected.GET("/users", middlewares.RoleMiddleware(enums.Admin, enums.Member), userHandler.GetAll)
			protected.GET("/perfil/:username", userHandler.Perfil)

			// Posts
			protected.POST("/posts", postHandler.CreatePost)
			protected.GET("/posts", postHandler.GetAllPosts)
			protected.GET("/cartelera", postHandler.Cartelera)
			protected.GET("/posts/:id", postHandler.GetPostById)
			protected.PUT("/posts", postHandler.UpdatePost)
			protected.DELETE("/posts/:id", postHandler.DeletePost)
			protected.POST("/posts/:id/vote", postHandler.PostVote)
			protected.GET("/users/:ownerId/posts", postHandler.GetPostsByUserId)

			// Comments
			protected.POST("/comments", commentHandler.CreateComment)
			protected.GET("/posts/:id/comments", commentHandler.GetCommentsByPostId)
			protected.PUT("/comments/:id", commentHandler.UpdateComment)
			protected.DELETE("/comments/:id", commentHandler.DeleteComment)
		}
	}
}
