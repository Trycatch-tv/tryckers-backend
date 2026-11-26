package handlers

import (
	"io"
	"net/http"
	"strings"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/services/storage"
	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	Storage storage.StorageService
}

// GetMedia recupera y transmite un objeto multimedia desde el driver de almacenamiento configurado (usado en Dev).
func (h *MediaHandler) GetMedia(c *gin.Context) {
	if h.Storage == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	rawKey := strings.TrimPrefix(c.Param("key"), "/")
	if !storage.IsManagedObjectKey(rawKey) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "clave de media inválida"})
		return
	}

	reader, contentType, err := h.Storage.GetObject(c.Request.Context(), rawKey)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	defer reader.Close()

	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, reader)
}
