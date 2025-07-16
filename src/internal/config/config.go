package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
}

func Load(env string) Config {

	var envFile string
	switch env {
	case "production":
		envFile = ".env.production"
	case "test":
		envFile = ".env.test"
	default:
		envFile = ".env"
	}

	rootPath := findProjectRoot()
	err := godotenv.Load(filepath.Join(rootPath, envFile))
	if err != nil {
		log.Fatalf("Error loading %s: %v", envFile, err)
	}

	fmt.Print("✅ load file  ", envFile)

	return Config{
		Port:              getEnv("PORT", "8080"),
		DBUrl:             getEnv("DATABASE_URL", ""),
		POSTGRES_DB:       getEnv("POSTGRES_DB", ""),
		POSTGRES_USER:     getEnv("POSTGRES_USER", ""),
		POSTGRES_PASSWORD: getEnv("POSTGRES_PASSWORD", ""),
		POSTGRES_HOST:     getEnv("POSTGRES_HOST", "localhost"),
		POSTGRES_PORT:     getEnv("POSTGRES_PORT", "5432"),
		JWT_SECRET:        getEnv("JWT_SECRET", ""),
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

func findProjectRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break // llegó a la raíz
		}
		dir = parent
	}
	log.Fatal("❌ No se encontró el archivo go.mod. ¿Estás ejecutando desde dentro de un módulo Go?")
	return ""
}
