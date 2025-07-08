package tests

import (
	"log"
	"os"
	"testing"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/config"
)

func TestMain(m *testing.M) {
	var db = config.InitGormDB(config.LoadTest())

	// Run tests
	code := m.Run()

	// After all tests: reset DB
	log.Println("🧹 Cleaning up test DB...")
	// todo: cada vez que se agregue una entidad o una nueva tabla se debe agregar a la query de limpiar la db de testing
	db.Exec("TRUNCATE users, posts, comments RESTART IDENTITY CASCADE;")

	os.Exit(code)
}
