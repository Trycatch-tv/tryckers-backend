package handlers

import (
	"net/http"
	"strings"

	"github.com/Trycatch-tv/tryckers-backend/src/internal/dtos"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/enums"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/services"
	"github.com/Trycatch-tv/tryckers-backend/src/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
// @Success      200  {array}   dtos.CreateUserDTO  "List of users"
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
// @Success      201  {object}  dtos.CreateUserDTO  "Created user"
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
// @Success      200  {object}  dtos.LoginResponse  "User information"
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
// @Success      200  {object}  dtos.LoginResponse  "User profile information"
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

// UploadAvatar godoc
// @Summary      Upload current user avatar
// @Description  Upload and replace the authenticated user's profile avatar
// @Tags         Profile
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Avatar image"
// @Success      200   {object}  models.User  "Updated user"
// @Failure      400   {object}  ErrorResponse  "Invalid file"
// @Failure      401   {object}  ErrorResponse  "Unauthorized"
// @Security     BearerAuth
// @Router       /users/me/avatar [post]
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	h.uploadMedia(c, "avatar")
}

// UploadBanner godoc
// @Summary      Upload current user banner
// @Description  Upload and replace the authenticated user's profile banner
// @Tags         Profile
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Banner image"
// @Success      200   {object}  models.User  "Updated user"
// @Failure      400   {object}  ErrorResponse  "Invalid file"
// @Failure      401   {object}  ErrorResponse  "Unauthorized"
// @Security     BearerAuth
// @Router       /users/me/banner [post]
func (h *UserHandler) UploadBanner(c *gin.Context) {
	h.uploadMedia(c, "banner")
}

// DeleteAvatar godoc
// @Summary      Remove current user avatar
// @Description  Remove the authenticated user's profile avatar
// @Tags         Profile
// @Produce      json
// @Success      200  {object}  models.User  "Updated user"
// @Failure      401  {object}  ErrorResponse  "Unauthorized"
// @Security     BearerAuth
// @Router       /users/me/avatar [delete]
func (h *UserHandler) DeleteAvatar(c *gin.Context) {
	h.deleteMedia(c, "avatar")
}

// DeleteBanner godoc
// @Summary      Remove current user banner
// @Description  Remove the authenticated user's profile banner
// @Tags         Profile
// @Produce      json
// @Success      200  {object}  models.User  "Updated user"
// @Failure      401  {object}  ErrorResponse  "Unauthorized"
// @Security     BearerAuth
// @Router       /users/me/banner [delete]
func (h *UserHandler) DeleteBanner(c *gin.Context) {
	h.deleteMedia(c, "banner")
}

func (h *UserHandler) uploadMedia(c *gin.Context, mediaType string) {
	userID, ok := h.currentUserID(c)
	if !ok {
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		HandleBadRequest(c, "archivo requerido")
		return
	}

	var updatedUser interface{}
	if mediaType == "avatar" {
		updatedUser, err = h.Service.UploadAvatar(userID, file)
	} else {
		updatedUser, err = h.Service.UploadBanner(userID, file)
	}

	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": updatedUser})
}

func (h *UserHandler) deleteMedia(c *gin.Context, mediaType string) {
	userID, ok := h.currentUserID(c)
	if !ok {
		return
	}

	var updatedUser interface{}
	var err error
	if mediaType == "avatar" {
		updatedUser, err = h.Service.RemoveAvatar(userID)
	} else {
		updatedUser, err = h.Service.RemoveBanner(userID)
	}

	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": updatedUser})
}

func (h *UserHandler) currentUserID(c *gin.Context) (uuid.UUID, bool) {
	rawUserID, exists := c.Get("userId")
	if !exists {
		HandleUnauthorized(c, "usuario no autenticado")
		return uuid.Nil, false
	}

	userIDValue, ok := rawUserID.(string)
	if !ok {
		HandleBadRequest(c, "usuario invalido")
		return uuid.Nil, false
	}

	userID, err := uuid.Parse(userIDValue)
	if err != nil {
		HandleBadRequest(c, "usuario invalido")
		return uuid.Nil, false
	}

	return userID, true
}
