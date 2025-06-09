package main

import (
	"log"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/api/routes"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/config"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db := config.InitGormDB(cfg)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	db.AutoMigrate(&models.User{})

	r := gin.Default()
	routes.SetupV1(r, db)

	log.Println("Server running on port", cfg.Port)
	r.Run(":" + cfg.Port)
}
