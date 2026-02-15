package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"

	apperrors "github.com/Trycatch-tv/tryckers-backend/src/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// APIError representa la estructura de respuesta de error para la API
type APIError struct {
	Error   string            `json:"error" example:"Error message"`
	Code    int               `json:"code,omitempty" example:"400"`
	Details string            `json:"details,omitempty" example:"Additional details"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// validationMessages mapea tags de validación a mensajes en español
var validationMessages = map[string]string{
	"required": "es requerido",
	"min":      "debe tener al menos %s caracteres",
	"max":      "debe tener máximo %s caracteres",
	"url":      "debe ser una URL válida",
	"email":    "debe ser un email válido",
	"oneof":    "debe ser uno de los valores permitidos: %s",
	"uuid":     "debe ser un UUID válido",
}

// getValidationMessage genera un mensaje de error amigable para un error de validación
func getValidationMessage(fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()
	param := fe.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("el campo '%s' es requerido", field)
	case "min":
		return fmt.Sprintf("el campo '%s' debe tener al menos %s caracteres", field, param)
	case "max":
		return fmt.Sprintf("el campo '%s' debe tener máximo %s caracteres", field, param)
	case "url":
		return fmt.Sprintf("el campo '%s' debe ser una URL válida", field)
	case "email":
		return fmt.Sprintf("el campo '%s' debe ser un email válido", field)
	case "oneof":
		return fmt.Sprintf("el campo '%s' debe ser uno de: %s", field, param)
	case "uuid":
		return fmt.Sprintf("el campo '%s' debe ser un UUID válido", field)
	default:
		return fmt.Sprintf("el campo '%s' no es válido", field)
	}
}

// HandleBindingError maneja errores de binding/validación de Gin y devuelve mensajes descriptivos
func HandleBindingError(c *gin.Context, err error) {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		fields := make(map[string]string)
		var firstError string

		for i, fe := range validationErrors {
			msg := getValidationMessage(fe)
			fields[fe.Field()] = msg
			if i == 0 {
				firstError = msg
			}
		}

		c.JSON(http.StatusBadRequest, APIError{
			Error:  firstError,
			Code:   http.StatusBadRequest,
			Fields: fields,
		})
		return
	}

	// Si no es un error de validación, devolver mensaje genérico
	c.JSON(http.StatusBadRequest, APIError{
		Error: "datos de entrada inválidos",
		Code:  http.StatusBadRequest,
	})
}

// HandleError maneja los errores y envía la respuesta HTTP apropiada
func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	// Log del error para debugging
	log.Printf("Error: %v", err)

	// Si es un AppError, usar su código y mensaje
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		c.JSON(appErr.Code, APIError{
			Error: appErr.Message,
			Code:  appErr.Code,
		})
		return
	}

	// Verificar si es un error de GORM (record not found)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, APIError{
			Error: "recurso no encontrado",
			Code:  http.StatusNotFound,
		})
		return
	}

	// Error genérico - devolver 500
	c.JSON(http.StatusInternalServerError, APIError{
		Error: "error interno del servidor",
		Code:  http.StatusInternalServerError,
	})
}

// ParseUUID valida y parsea un string a UUID
// Retorna el UUID parseado y un error si el formato es inválido
func ParseUUID(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.Nil, apperrors.ErrInvalidID
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, apperrors.NewBadRequestWithError("formato de UUID inválido", err)
	}

	if parsedID == uuid.Nil {
		return uuid.Nil, apperrors.ErrInvalidID
	}

	return parsedID, nil
}

// HandleBadRequest envía una respuesta de error 400
func HandleBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, APIError{
		Error: message,
		Code:  http.StatusBadRequest,
	})
}

// HandleNotFound envía una respuesta de error 404
func HandleNotFound(c *gin.Context, resource string) {
	c.JSON(http.StatusNotFound, APIError{
		Error: resource + " no encontrado",
		Code:  http.StatusNotFound,
	})
}

// HandleUnauthorized envía una respuesta de error 401
func HandleUnauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, APIError{
		Error: message,
		Code:  http.StatusUnauthorized,
	})
}

// HandleInternalError envía una respuesta de error 500
func HandleInternalError(c *gin.Context, err error) {
	log.Printf("Internal error: %v", err)
	c.JSON(http.StatusInternalServerError, APIError{
		Error: "error interno del servidor",
		Code:  http.StatusInternalServerError,
	})
}

// HandleValidationError envía una respuesta de error 422
func HandleValidationError(c *gin.Context, message string) {
	c.JSON(http.StatusUnprocessableEntity, APIError{
		Error: message,
		Code:  http.StatusUnprocessableEntity,
	})
}
