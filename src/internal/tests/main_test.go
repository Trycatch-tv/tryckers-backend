package tests

import (
	"log"
	"os"
	"testing"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/config"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
)

func TestMain(m *testing.M) {
	env := "test"
	Testdb = config.InitGormDB(config.Load(env))

	Testdb.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	Testdb.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})
	// Run tests
	code := m.Run()

	// After all tests: reset DB
	log.Println("🧹 Cleaning up test DB...")
	// todo: cada vez que se agregue una entidad o una nueva tabla se debe agregar a la query de limpiar la db de testing
	Testdb.Exec("TRUNCATE users, posts, comments RESTART IDENTITY CASCADE;")

	os.Exit(code)
}
