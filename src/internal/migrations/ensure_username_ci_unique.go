package migrations

import (
	"log"

	"gorm.io/gorm"
)

// EnsureUsernameCIUniqueIndex crea un índice único case-insensitive (y sin espacios al borde)
// sobre users.username para impedir duplicados lógicos (ej. "Deiby7h322" vs "deiby7h322"),
// incluyendo registros históricos y solicitudes concurrentes.
func EnsureUsernameCIUniqueIndex(db *gorm.DB) error {
	stmt := `CREATE UNIQUE INDEX IF NOT EXISTS idx_users_username_lower_unique ON users (LOWER(TRIM(username)))`
	if err := db.Exec(stmt).Error; err != nil {
		return err
	}
	log.Println("[migration] Ensured case-insensitive unique index idx_users_username_lower_unique on users")
	return nil
}
