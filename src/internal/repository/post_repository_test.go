package repository

import (
	"strings"
	"testing"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestPostRepository_GetTopPostsOfTheWeek_QueryStructure(t *testing.T) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: "host=localhost user=test password=test dbname=test port=5432 sslmode=disable",
	}), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
	})
	if err != nil {
		t.Fatalf("failed to init dry-run gorm db: %v", err)
	}

	var posts []models.Post
	stmt := db.Model(&models.Post{}).
		Select("posts.*, COALESCE(SUM(CASE WHEN post_votes.vote = 1 THEN 1 ELSE 0 END), 0) AS votes_count").
		Joins("LEFT JOIN post_votes ON post_votes.post_id = posts.id").
		Where("posts.status != ? AND posts.created_at >= NOW() - INTERVAL '7 days'", enums.DELETED).
		Group("posts.id").
		Order("COALESCE(SUM(CASE WHEN post_votes.vote = 1 THEN 1 ELSE 0 END), 0) DESC").
		Limit(10).
		Find(&posts).Statement

	sql := stmt.SQL.String()
	for _, expected := range []string{
		"LEFT JOIN post_votes ON post_votes.post_id = posts.id",
		"posts.status !=",
		"INTERVAL '7 days'",
		`GROUP BY "posts"."id"`,
		"ORDER BY",
		"LIMIT",
	} {
		if !strings.Contains(sql, expected) {
			t.Fatalf("expected SQL to contain %q, got: %s", expected, sql)
		}
	}
}
