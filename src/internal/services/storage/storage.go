package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

const (
	AvatarKind = "avatars"
	BannerKind = "banners"
	PostKind   = "posts"
)

var allowedImageTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

type SaveOptions struct {
	Kind     string
	OwnerID  string
	MaxBytes int64
}

type ValidatedImage struct {
	Data        []byte
	ContentType string
	Extension   string
}

// StorageService desacopla las operaciones multimedia del proveedor de almacenamiento (S3 o Local).
type StorageService interface {
	// SaveImage valida y persiste la imagen, retornando su object_key canónico (ej. avatars/{ownerId}/{uuid}.webp).
	SaveImage(ctx context.Context, fileHeader *multipart.FileHeader, options SaveOptions) (string, error)
	// GetObject recupera el stream y Content-Type de un objeto dado su object_key (usado en ruta /media en Dev).
	GetObject(ctx context.Context, objectKey string) (io.ReadCloser, string, error)
	// DeleteObject elimina un objeto dado su object_key canónico.
	DeleteObject(ctx context.Context, objectKey string) error
	// ResolveURL construye la URL de acceso para un object_key válido según el entorno.
	// Retorna "" si la referencia está vacía o corresponde a una URL/ruta antigua no canónica.
	ResolveURL(objectKey string) string
}

// IsManagedObjectKey verifica si una cadena corresponde a una clave canónica administrada por Tryckers.
func IsManagedObjectKey(value string) bool {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || strings.Contains(trimmed, "..") || strings.HasPrefix(trimmed, "/") {
		return false
	}
	return strings.HasPrefix(trimmed, AvatarKind+"/") ||
		strings.HasPrefix(trimmed, BannerKind+"/") ||
		strings.HasPrefix(trimmed, PostKind+"/")
}

// GenerateObjectKey construye una clave única canónica: {kind}/{ownerId}/{uuid}{ext}.
func GenerateObjectKey(kind, ownerID, extension string) string {
	safeOwner := strings.NewReplacer("/", "-", "\\", "-", " ", "-", "..", "-").Replace(strings.TrimSpace(ownerID))
	if safeOwner == "" {
		safeOwner = "unknown"
	}
	return fmt.Sprintf("%s/%s/%s%s", kind, safeOwner, uuid.NewString(), extension)
}

// ResolveMediaURL convierte un object_key canónico en una URL pública o relativa según baseURL.
// No ofrece fallback para URLs o rutas antiguas (/uploads/... o http://...).
func ResolveMediaURL(baseURL, objectKey string) string {
	if !IsManagedObjectKey(objectKey) {
		return ""
	}
	trimmedBase := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if trimmedBase == "" {
		return "/media/" + objectKey
	}
	return trimmedBase + "/" + objectKey
}

// ValidateAndReadImage centraliza la validación de tipo de media, tamaño máximo y MIME real.
func ValidateAndReadImage(fileHeader *multipart.FileHeader, options SaveOptions) (*ValidatedImage, error) {
	if fileHeader == nil {
		return nil, fmt.Errorf("archivo requerido")
	}

	if options.Kind != AvatarKind && options.Kind != BannerKind && options.Kind != PostKind {
		return nil, fmt.Errorf("tipo de media invalido")
	}

	if fileHeader.Size <= 0 {
		return nil, fmt.Errorf("archivo vacio")
	}

	if options.MaxBytes > 0 && fileHeader.Size > options.MaxBytes {
		return nil, fmt.Errorf("archivo supera el tamano maximo permitido")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir el archivo: %w", err)
	}
	defer file.Close()

	limited := io.LimitReader(file, options.MaxBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer el archivo: %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("archivo vacio")
	}
	if options.MaxBytes > 0 && int64(len(data)) > options.MaxBytes {
		return nil, fmt.Errorf("archivo supera el tamano maximo permitido")
	}

	sniffLen := len(data)
	if sniffLen > 512 {
		sniffLen = 512
	}
	mimeType := http.DetectContentType(data[:sniffLen])
	extension, ok := allowedImageTypes[mimeType]
	if !ok {
		return nil, fmt.Errorf("tipo de archivo no permitido")
	}

	return &ValidatedImage{
		Data:        data,
		ContentType: mimeType,
		Extension:   extension,
	}, nil
}

// NewBytesReader retorna un io.ReadSeeker reutilizable para cargas a S3.
func NewBytesReader(data []byte) *bytes.Reader {
	return bytes.NewReader(data)
}
