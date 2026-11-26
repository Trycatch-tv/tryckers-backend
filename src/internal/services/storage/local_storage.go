package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// LocalStorage se conserva únicamente como fallback para pruebas unitarias y entornos sin conexión.
type LocalStorage struct {
	RootDir       string
	PublicBaseURL string
}

func NewLocalStorage(rootDir, publicBaseURL string) *LocalStorage {
	return &LocalStorage{
		RootDir:       rootDir,
		PublicBaseURL: strings.TrimRight(strings.TrimSpace(publicBaseURL), "/"),
	}
}

func (s *LocalStorage) SaveImage(_ context.Context, fileHeader *multipart.FileHeader, options SaveOptions) (string, error) {
	img, err := ValidateAndReadImage(fileHeader, options)
	if err != nil {
		return "", err
	}

	objectKey := GenerateObjectKey(options.Kind, options.OwnerID, img.Extension)
	targetPath := filepath.Join(s.RootDir, filepath.FromSlash(objectKey))

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return "", fmt.Errorf("no se pudo crear el directorio de media: %w", err)
	}

	if err := os.WriteFile(targetPath, img.Data, 0644); err != nil {
		return "", fmt.Errorf("no se pudo escribir el archivo: %w", err)
	}

	return objectKey, nil
}

func (s *LocalStorage) GetObject(_ context.Context, objectKey string) (io.ReadCloser, string, error) {
	if !IsManagedObjectKey(objectKey) {
		return nil, "", fmt.Errorf("object_key invalido")
	}

	targetPath := filepath.Join(s.RootDir, filepath.FromSlash(objectKey))
	file, err := os.Open(targetPath)
	if err != nil {
		return nil, "", err
	}

	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	_, _ = file.Seek(0, io.SeekStart)
	contentType := http.DetectContentType(buf[:n])

	return file, contentType, nil
}

func (s *LocalStorage) DeleteObject(_ context.Context, objectKey string) error {
	if !IsManagedObjectKey(objectKey) {
		return nil
	}

	cleanPath := filepath.Clean(filepath.FromSlash(objectKey))
	if strings.HasPrefix(cleanPath, "..") || filepath.IsAbs(cleanPath) {
		return nil
	}

	err := os.Remove(filepath.Join(s.RootDir, cleanPath))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (s *LocalStorage) ResolveURL(objectKey string) string {
	return ResolveMediaURL(s.PublicBaseURL, objectKey)
}
