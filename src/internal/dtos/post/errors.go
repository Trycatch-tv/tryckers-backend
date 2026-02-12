package dtos

import "errors"

// Errores de validación para DTOs de posts
var (
	// ErrVideoPostRequiresMedia indica que un post de tipo video requiere una URL de media
	ErrVideoPostRequiresMedia = errors.New("los posts de tipo video requieren una URL de media")

	// ErrInvalidPostType indica que el tipo de post no es válido
	ErrInvalidPostType = errors.New("el tipo de post no es válido")

	// ErrInvalidPostStatus indica que el estado del post no es válido
	ErrInvalidPostStatus = errors.New("el estado del post no es válido")

	// ErrTitleTooShort indica que el título es demasiado corto
	ErrTitleTooShort = errors.New("el título debe tener al menos 3 caracteres")

	// ErrTitleTooLong indica que el título es demasiado largo
	ErrTitleTooLong = errors.New("el título no puede exceder los 500 caracteres")

	// ErrContentTooShort indica que el contenido es demasiado corto
	ErrContentTooShort = errors.New("el contenido debe tener al menos 10 caracteres")

	// ErrContentTooLong indica que el contenido es demasiado largo
	ErrContentTooLong = errors.New("el contenido no puede exceder los 10000 caracteres")

	// ErrTooManyTags indica que hay demasiados tags
	ErrTooManyTags = errors.New("no se pueden agregar más de 10 tags")

	// ErrTagTooShort indica que un tag es demasiado corto
	ErrTagTooShort = errors.New("cada tag debe tener al menos 2 caracteres")

	// ErrTagTooLong indica que un tag es demasiado largo
	ErrTagTooLong = errors.New("cada tag no puede exceder los 50 caracteres")

	// ErrInvalidImageURL indica que la URL de la imagen no es válida
	ErrInvalidImageURL = errors.New("la URL de la imagen no es válida")

	// ErrInvalidMediaURL indica que la URL del media no es válida
	ErrInvalidMediaURL = errors.New("la URL del media no es válida")

	// ErrNoChangesProvided indica que no se proporcionaron cambios en la actualización
	ErrNoChangesProvided = errors.New("no se proporcionaron cambios para actualizar")
)
