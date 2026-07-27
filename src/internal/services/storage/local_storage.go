package storage

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	AvatarKind = "avatars"
	BannerKind = "banners"
)

var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

type LocalStorage struct {
	RootDir   string
	PublicURL string
}

type SaveOptions struct {
	Kind     string
	UserID   string
	MaxBytes int64
}

func NewLocalStorage(rootDir, publicURL string) *LocalStorage {
	return &LocalStorage{
		RootDir:   rootDir,
		PublicURL: strings.TrimRight(publicURL, "/"),
	}
}

func (s *LocalStorage) SaveImage(fileHeader *multipart.FileHeader, options SaveOptions) (string, error) {
	if fileHeader == nil {
		return "", fmt.Errorf("archivo requerido")
	}

	if options.Kind != AvatarKind && options.Kind != BannerKind {
		return "", fmt.Errorf("tipo de media invalido")
	}

	if fileHeader.Size <= 0 {
		return "", fmt.Errorf("archivo vacio")
	}

	if fileHeader.Size > options.MaxBytes {
		return "", fmt.Errorf("archivo supera el tamano maximo permitido")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("no se pudo abrir el archivo: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	bytesRead, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("no se pudo leer el archivo: %w", err)
	}

	mimeType := http.DetectContentType(buffer[:bytesRead])
	extension, ok := allowedImageTypes[mimeType]
	if !ok {
		return "", fmt.Errorf("tipo de archivo no permitido")
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("no se pudo preparar el archivo: %w", err)
	}

	targetDir := filepath.Join(s.RootDir, options.Kind)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("no se pudo crear el directorio de media: %w", err)
	}

	fileName := s.safeFileName(options.Kind, options.UserID, extension)
	targetPath := filepath.Join(targetDir, fileName)

	target, err := os.Create(targetPath)
	if err != nil {
		return "", fmt.Errorf("no se pudo guardar el archivo: %w", err)
	}
	defer target.Close()

	if _, err := io.Copy(target, file); err != nil {
		_ = os.Remove(targetPath)
		return "", fmt.Errorf("no se pudo escribir el archivo: %w", err)
	}

	return fmt.Sprintf("%s/%s/%s", s.PublicURL, options.Kind, fileName), nil
}

func (s *LocalStorage) DeleteByPublicURL(publicURL string) {
	if publicURL == "" || !strings.HasPrefix(publicURL, s.PublicURL+"/") {
		return
	}

	relativePath := strings.TrimPrefix(publicURL, s.PublicURL+"/")
	cleanPath := filepath.Clean(relativePath)
	if strings.HasPrefix(cleanPath, "..") || filepath.IsAbs(cleanPath) {
		return
	}

	_ = os.Remove(filepath.Join(s.RootDir, cleanPath))
}

func (s *LocalStorage) safeFileName(kind, userID, extension string) string {
	randomBytes := make([]byte, 8)
	if _, err := rand.Read(randomBytes); err != nil {
		copy(randomBytes, []byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	}

	safeUserID := strings.NewReplacer("/", "-", "\\", "-", " ", "-").Replace(userID)
	return fmt.Sprintf("%s-%s-%d-%s%s", kind, safeUserID, time.Now().UnixNano(), hex.EncodeToString(randomBytes), extension)
}
