package dtos

type LoginUser struct {
	Email string `json:"email" binding:"required,email"` // ¿Es correcto usar el email como usuario?
	// Si el email se va a publicar como
	// una forma de contacto, lo más seguro es que utilicen el mismo email.

	Password string `json:"password" binding:"required"`
}
