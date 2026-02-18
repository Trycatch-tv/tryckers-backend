package handlers

import (
	"net/http"
	"strings"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/dtos"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/services"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *services.UserService
}

// ErrorResponse representa una respuesta de error estándar
type ErrorResponse struct {
	Error string `json:"error" example:"Error message"`
} // @name ErrorResponse

// GetAll godoc
// @Summary      Get all users
// @Description  Retrieve a list of all registered users in the system
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.User  "List of users"
// @Failure      500  {object}  ErrorResponse  "Internal server error"
// @Security     BearerAuth
// @Router       /users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		HandleError(c, err)
		return
	}
	c.JSON(http.StatusOK, users)
}

// CreateUser godoc
// @Summary      Register a new user
// @Description  Register a new user in the system with the provided information
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user  body      dtos.CreateUserDTO  true  "User registration information"
// @Success      201  {object}  models.User  "Created user"
// @Failure      400  {object}  ErrorResponse  "Invalid input or country"
// @Failure      500  {object}  ErrorResponse  "Internal server error"
// @Router       /register [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var newUser dtos.CreateUserDTO

	if err := c.ShouldBindJSON(&newUser); err != nil {
		HandleBindingError(c, err)
		return
	}

	if !enums.IsValidCountry(string(newUser.Country)) {
		HandleBadRequest(c, "país inválido")
		return
	}

	if !h.Service.IsEmailRegistered(newUser.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the email is already registered"})
		return
	}

	userCreated, err := h.Service.CreateUser(&newUser)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, userCreated)
}

// Login godoc
// @Summary      User login
// @Description  Authenticate a user and return their information with an access token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      dtos.LoginUser  true  "Login credentials"
// @Success      200  {object}  object{user=models.User}  "User information"
// @Failure      400  {object}  ErrorResponse  "Invalid credentials format"
// @Failure      500  {object}  ErrorResponse  "Authentication failed or internal server error"
// @Router       /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var user dtos.LoginUser

	if err := c.ShouldBindJSON(&user); err != nil {
		HandleBindingError(c, err)
		return
	}

	loginResponse, err := h.Service.Login(&user)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": loginResponse})
}

// Profile godoc
// @Summary      Get user profile by username
// @Description  Retrieve detailed profile information for a specific user by their username
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Param        email  path      string  true  "email of the user"
// @Success      200  {object}  object{user=models.User}  "User profile information"
// @Failure      404  {object}  ErrorResponse  "User not found"
// @Security     BearerAuth
// @Router       /profile/{username} [get]
func (h *UserHandler) Perfil(c *gin.Context) {
	// username es en toLowerCase
	username := strings.ToLower(c.Param("username"))

	if username == "" {
		HandleBadRequest(c, "username requerido")
		return
	}

	userPerfil, err := h.Service.Perfil(username)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userPerfil})
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Generate new access and refresh tokens using a valid refresh token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        refresh_token  body      dtos.RefreshTokenRequest  true  "Refresh token"
// @Success      200  {object}  dtos.RefreshTokenResponse  "New tokens"
// @Failure      400  {object}  ErrorResponse  "Invalid request format"
// @Failure      401  {object}  ErrorResponse  "Invalid or expired refresh token"
// @Router       /refresh-token [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req dtos.RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		HandleBadRequest(c, "refresh_token es requerido")
		return
	}

	accessToken, refreshToken, err := utils.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error(), "code": 401})
		return
	}

	c.JSON(http.StatusOK, dtos.RefreshTokenResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
