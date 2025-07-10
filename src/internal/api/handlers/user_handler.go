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

func (h *UserHandler) GetAll(c *gin.Context) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

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

	userCreated, err := h.Service.CreateUser(&newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, userCreated)
}

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

func (h *UserHandler) Perfil(c *gin.Context) {
	email := c.Param("email")

	userPerfil, err := h.Service.Perfil(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": userPerfil})
}
