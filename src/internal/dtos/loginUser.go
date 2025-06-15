package dtos

// LoginUser represents the login request payload
type LoginUser struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"mypassword123"`
} // @name LoginUser
