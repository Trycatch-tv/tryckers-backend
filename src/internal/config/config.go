package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Port              string
	DBUrl             string
	POSTGRES_DB       string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
	POSTGRES_HOST     string
	POSTGRES_PORT     string
	JWT_SECRET        string
	FRONTEND_URL      string
}

func Load() Config {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		Port:              getEnv("PORT", "8080"),
		DBUrl:             getEnv("DATABASE_URL", ""),
		POSTGRES_DB:       getEnv("POSTGRES_DB", ""),
		POSTGRES_USER:     getEnv("POSTGRES_USER", ""),
		POSTGRES_PASSWORD: getEnv("POSTGRES_PASSWORD", ""),
		POSTGRES_HOST:     getEnv("POSTGRES_HOST", "localhost"),
		POSTGRES_PORT:     getEnv("POSTGRES_PORT", "5432"),
		JWT_SECRET:        getEnv("JWT_SECRET", ""),
		FRONTEND_URL:      getEnv("FRONTEND_URL", "http://localhost:4200"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func InitGormDB(cfg Config) *gorm.DB {
	dbUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.POSTGRES_USER,
		cfg.POSTGRES_PASSWORD,
		cfg.POSTGRES_HOST,
		cfg.POSTGRES_PORT,
		cfg.POSTGRES_DB,
	)
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}
	return db
}
