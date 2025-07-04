package routes

import (
	"github.com/Trycatch-tv/tryckers-backend/src/internal/api/handlers"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/repository"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupV1(r *gin.Engine, db *gorm.DB) {
	userRepo := &repository.UserRepository{DB: db}
	userService := &services.UserService{Repo: userRepo}
	userHandler := &handlers.UserHandler{Service: userService}
	commentRepo := &repository.CommentRepository{DB: db}
	commentService := &services.CommentService{Repo: commentRepo}
	commentHandler := &handlers.CommentHandler{Service: commentService}
	postRepo := &repository.PostRepository{DB: db}
	postService := &services.PostService{Repo: postRepo}
	postHandler := &handlers.PostHandler{Service: postService}
	api := r.Group("/api/v1")
	{
		api.GET("/users", userHandler.GetAll)
		api.POST("/register", userHandler.CreateUser)
		api.POST("/login", userHandler.Login)
		api.POST("/comments", commentHandler.CreateComment)
		api.GET("/posts/:id/comments", commentHandler.GetCommentsByPostId)
		api.PUT("/comments/:id", commentHandler.UpdateComment)
		api.DELETE("/comments/:id", commentHandler.DeleteComment)
		api.POST("/posts", postHandler.CreatePost)
		api.GET("/posts", postHandler.GetAllPosts)
		api.GET("/posts/:id", postHandler.GetPostById)
		api.PUT("/posts", postHandler.UpdatePost)
		api.DELETE("/posts/:id", postHandler.DeletePost)
	}
}
