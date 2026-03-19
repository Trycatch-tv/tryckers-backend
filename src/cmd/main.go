package main

import (
	"fmt"
	"log"

	_ "github.com/Trycatch-tv/tryckers-backend/docs"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/api/routes"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/config"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Tryckers API
// @version         1.0
// @description     A modern API for the Tryckers platform, providing user management and authentication services.
// @description     This API allows you to manage users, handle authentication, and access user profiles.

// @contact.name   API Support
// @contact.url    https://tryckers.com/support
// @contact.email  support@tryckers.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @type apiKey
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer abcde12345"

// @tag.name Users
// @tag.description User management operations

// @tag.name Auth
// @tag.description Authentication operations

// @tag.name Profile
// @tag.description User profile operations

// @tag.name Posts
// @tag.description Post management operations

// @tag.name Comments
// @tag.description Comment management operations

func main() {
	log.Println("Iniciando backend...")
	cfg := config.Load()
	db := config.InitGormDB(cfg)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Enable UUID extension
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{}, &models.PostVote{})

	// Initialize Gin with default middleware
	r := gin.Default()

	// configuración de cors
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FRONTEND_URL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	// Servir archivos estáticos de la carpeta docs
	r.Static("/docs", "./docs")

	swaggerURL := fmt.Sprintf("http://localhost:%s/docs/swagger.json", cfg.Port)

	// Swagger documentation endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL(swaggerURL),
		ginSwagger.DefaultModelsExpandDepth(-1),
	))

	// Setup API routes
	routes.SetupV1(r, db)

	log.Printf("Server running on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
