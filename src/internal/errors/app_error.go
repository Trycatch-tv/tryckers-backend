package errors

import (
	"fmt"
	"net/http"
)

// AppError representa un error de aplicación con código HTTP y mensaje
type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"error"`
	Err     error  `json:"-"`
}

// Error implementa la interfaz error
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap permite usar errors.Is y errors.As
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError crea un nuevo AppError con código, mensaje y error original
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewBadRequest crea un error 400 Bad Request
func NewBadRequest(message string) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}

// NewBadRequestWithError crea un error 400 con error original
func NewBadRequestWithError(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusBadRequest,
		Message: message,
		Err:     err,
	}
}

// NewNotFound crea un error 404 Not Found
func NewNotFound(resource string) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%s no encontrado", resource),
	}
}

// NewNotFoundWithError crea un error 404 con error original
func NewNotFoundWithError(resource string, err error) *AppError {
	return &AppError{
		Code:    http.StatusNotFound,
		Message: fmt.Sprintf("%s no encontrado", resource),
		Err:     err,
	}
}

// NewUnauthorized crea un error 401 Unauthorized
func NewUnauthorized(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
}

// NewForbidden crea un error 403 Forbidden
func NewForbidden(message string) *AppError {
	return &AppError{
		Code:    http.StatusForbidden,
		Message: message,
	}
}

// NewInternalError crea un error 500 Internal Server Error
func NewInternalError(message string, err error) *AppError {
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}

// NewConflict crea un error 409 Conflict
func NewConflict(message string) *AppError {
	return &AppError{
		Code:    http.StatusConflict,
		Message: message,
	}
}

// NewValidationError crea un error 422 Unprocessable Entity
func NewValidationError(message string) *AppError {
	return &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: message,
	}
}

// Wrap envuelve un error existente con contexto adicional manteniendo el código HTTP
func Wrap(err error, message string) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return &AppError{
			Code:    appErr.Code,
			Message: message,
			Err:     appErr,
		}
	}
	return &AppError{
		Code:    http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}
