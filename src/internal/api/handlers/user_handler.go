package handlers

import (
	"net/http"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/dtos"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/services"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !enums.IsValidCountry(string(newUser.Country)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid country"})
		return
	}

	if !h.Service.IsvalidEmail(newUser.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the email is already registered"})
		return
	}

	userCreated, err := h.Service.CreateUser(&newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loginResponse, err := h.Service.Login(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Router       /profile/{email} [get]
func (h *UserHandler) Perfil(c *gin.Context) {
	email := c.Param("email")

	userPerfil, err := h.Service.Perfil(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userPerfil})
}
