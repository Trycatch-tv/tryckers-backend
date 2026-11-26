package repository

import (
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestUserRepository_FindByUsername_IsCaseInsensitive verifica que la búsqueda de
// usuario por username no distinga mayúsculas de minúsculas. Un usuario persistido
// con mayúsculas (p. ej. "Deiby7h322") debe poder encontrarse igual, por lo que el
// SQL que genera FindByUsername debe comparar normalizando con LOWER().
func TestUserRepository_FindByUsername_IsCaseInsensitive(t *testing.T) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "host=localhost user=test password=test dbname=test port=5432 sslmode=disable",
	}), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
	})
	if err != nil {
		t.Fatalf("failed to init dry-run gorm db: %v", err)
	}

	// Captura el SQL realmente generado por FindByUsername usando un callback de query.
	var capturedSQL string
	if err := db.Callback().Query().After("gorm:query").Register("capture_sql_test", func(tx *gorm.DB) {
		capturedSQL = tx.Statement.SQL.String()
	}); err != nil {
		t.Fatalf("failed to register query callback: %v", err)
	}

	repo := &UserRepository{DB: db}
	_, _ = repo.FindByUsername("Deiby7h322")

	sql := strings.ToLower(capturedSQL)
	if sql == "" {
		t.Fatalf("expected FindByUsername to build a query, captured empty SQL")
	}
	if !strings.Contains(sql, "lower(username) = lower(") {
		t.Fatalf("expected FindByUsername to compare with LOWER on both sides, got: %s", sql)
	}
}
