package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/brando1998/docubot-api/models"
)

// CreateClient godoc
// @Summary Crear un nuevo usuario
// @Description Crea un nuevo usuario en el sistema
// @Tags usuarios
// @Accept json
// @Produce json
// @Param user body map[string]string true "Datos del usuario"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /users [post]

func CreateClient(c *gin.Context) {
	var user models.Client
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}
	if err := clientRepo.CreateClient(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear usuario"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Usuario creado", "user": user})
}

// GetProfile godoc
// @Summary Obtener perfil del usuario
// @Description Retorna el perfil del usuario autenticado
// @Tags usuarios
// @Produce json
// @Success 200 {object} map[string]string
// @Router /users/profile [get]

func GetClientByID(c *gin.Context) {
	// Simula obtener el ID del usuario autenticado
	user, err := clientRepo.GetClientByID(1)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetOrCreateClient godoc
// @Summary Obtener o crear un usuario por teléfono
// @Description Busca un usuario por teléfono. Si no existe, lo crea.
// @Tags usuarios
// @Accept json
// @Produce json
// @Param data body map[string]string true "Datos del usuario"
// @Success 200 {object} models.Client
// @Failure 400 {object} map[string]string
// @Router /users/get-or-create [post]
func GetOrCreateClient(c *gin.Context) {
	var input struct {
		Phone string `json:"phone" binding:"required"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil || input.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Número de teléfono requerido"})
		return
	}

	user, err := clientRepo.GetOrCreateClient(input.Phone, input.Name, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar usuario"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetClientByPhone godoc
// @Summary Obtener usuario por teléfono
// @Description Busca un usuario por su número de teléfono
// @Tags usuarios
// @Produce json
// @Param phone path string true "Número de teléfono"
// @Success 200 {object} models.Client
// @Failure 404 {object} map[string]string
// @Router /users/{phone} [get]
func GetClientByPhone(c *gin.Context) {
	phone := c.Param("phone")
	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Número de teléfono requerido"})
		return
	}

	user, err := clientRepo.GetClientByPhone(phone)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, user)
}
