package main

import (
	"flag"
	"log"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/api/routes"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/config"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"github.com/gin-gonic/gin"
)

func main() {

	env := flag.String("app_env", "development", "Application environment (e.g. development, test, production)")
	flag.Parse()
	envString := *env

	cfg := config.Load(envString)
	db := config.InitGormDB(cfg)
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})

	r := gin.Default()
	routes.SetupV1(r, db)

	log.Println("Server running on port", cfg.Port)
	r.Run(":" + cfg.Port)
}
