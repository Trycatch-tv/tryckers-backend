package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/config"
)

// NewStorageFromConfig selecciona el driver de almacenamiento configurado (s3 o local).
func NewStorageFromConfig(ctx context.Context, cfg config.Config) (StorageService, error) {
	driver := strings.ToLower(strings.TrimSpace(cfg.StorageDriver))
	switch driver {
	case "s3":
		return NewS3Storage(ctx, cfg.AWSRegion, cfg.S3Bucket, cfg.MediaPublicBaseURL)
	case "", "local":
		baseURL := cfg.MediaPublicBaseURL
		if strings.TrimSpace(baseURL) == "" {
			baseURL = "/media"
		}
		return NewLocalStorage("uploads", baseURL), nil
	default:
		return nil, fmt.Errorf("STORAGE_DRIVER no soportado: %s", cfg.StorageDriver)
	}
}
