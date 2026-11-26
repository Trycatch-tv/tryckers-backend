package migrations

import (
	"log"

	"gorm.io/gorm"
)

// CleanupLegacyMediaReferences limpia de forma segura e idempotente las referencias
// multimedia antiguas (URLs o rutas /uploads/...) que no son object_key canónicos de S3.
// No modifica enlaces externos (github_url, linkedin_url, pitch_video, posts.media_url).
func CleanupLegacyMediaReferences(db *gorm.DB) error {
	statements := []string{
		`UPDATE users SET avatar_url = '' WHERE avatar_url IS NOT NULL AND avatar_url <> '' AND avatar_url NOT LIKE 'avatars/%'`,
		`UPDATE users SET profile_picture = '' WHERE profile_picture IS NOT NULL AND profile_picture <> '' AND profile_picture NOT LIKE 'avatars/%'`,
		`UPDATE users SET banner_url = '' WHERE banner_url IS NOT NULL AND banner_url <> '' AND banner_url NOT LIKE 'banners/%'`,
		`UPDATE posts SET image = '' WHERE image IS NOT NULL AND image <> '' AND image NOT LIKE 'posts/%'`,
		`UPDATE comments SET image = '' WHERE image IS NOT NULL AND image <> '' AND image NOT LIKE 'posts/%' AND image NOT LIKE 'comments/%'`,
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, stmt := range statements {
			res := tx.Exec(stmt)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected > 0 {
				log.Printf("[migration] Cleaned %d legacy media reference(s): %s", res.RowsAffected, stmt)
			}
		}
		return nil
	})
}
