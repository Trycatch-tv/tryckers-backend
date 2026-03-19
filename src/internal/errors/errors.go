package errors

import "errors"

// Errores predefinidos para uso común en el proyecto
var (
	// Errores de autenticación
	ErrInvalidCredentials = NewUnauthorized("credenciales incorrectas")
	ErrTokenInvalid       = NewUnauthorized("token inválido")
	ErrTokenExpired       = NewUnauthorized("token expirado")
	ErrUnauthorized       = NewUnauthorized("no autorizado")

	// Errores de validación
	ErrInvalidID            = NewBadRequest("ID inválido")
	ErrInvalidUUID          = NewBadRequest("UUID inválido")
	ErrInvalidInput         = NewBadRequest("datos de entrada inválidos")
	ErrMissingRequiredField = NewBadRequest("campo requerido faltante")

	// Errores de recursos
	ErrUserNotFound     = NewNotFound("usuario")
	ErrPostNotFound     = NewNotFound("post")
	ErrCommentNotFound  = NewNotFound("comentario")
	ErrResourceNotFound = NewNotFound("recurso")

	// Errores de duplicados
	ErrDuplicateUsername = NewConflict("el nombre de usuario ya está en uso")
	ErrDuplicateEmail    = NewConflict("el email ya está registrado")

	// Errores de permisos
	ErrForbidden     = NewForbidden("no tiene permisos para esta acción")
	ErrNotAuthorized = NewForbidden("usuario no autorizado")

	// Errores internos
	ErrInternalServer = NewInternalError("error interno del servidor", nil)
	ErrDatabaseError  = NewInternalError("error de base de datos", nil)
)

// IsNotFound verifica si un error es de tipo NotFound
func IsNotFound(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == 404
	}
	return false
}

// IsBadRequest verifica si un error es de tipo BadRequest
func IsBadRequest(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == 400
	}
	return false
}

// IsUnauthorized verifica si un error es de tipo Unauthorized
func IsUnauthorized(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == 401
	}
	return false
}

// IsForbidden verifica si un error es de tipo Forbidden
func IsForbidden(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == 403
	}
	return false
}

// IsInternalError verifica si un error es de tipo InternalError
func IsInternalError(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == 500
	}
	return false
}
