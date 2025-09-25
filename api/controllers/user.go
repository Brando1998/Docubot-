package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/brando1998/docubot-api/models"
)

// CreateClient godoc
// @Summary Crear un nuevo usuario
// @Description Crea un nuevo usuario en el sistema
// @Tags usuarios
// @Accept json
// @Produce json
// @Param user body models.Client true "Datos del usuario"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/v1/users [post]
func CreateClient(c *gin.Context) {
	var user models.Client
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	if err := clientRepo.CreateClient(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear usuario", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Usuario creado exitosamente",
		"user":    user,
	})
}

// GetClientByID godoc
// @Summary Obtener usuario por ID
// @Description Retorna un usuario específico por su ID
// @Tags usuarios
// @Produce json
// @Param id path int true "ID del usuario"
// @Success 200 {object} models.Client
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/users/id/{id} [get]
func GetClientByID(c *gin.Context) {
	// Obtener ID desde parámetros de URL
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de usuario requerido"})
		return
	}

	// Convertir string a uint
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Buscar usuario
	user, err := clientRepo.GetClientByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
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
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/users/phone/{phone} [get]
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

// GetOrCreateClient godoc
// @Summary Obtener o crear un usuario por teléfono
// @Description Busca un usuario por teléfono. Si no existe, lo crea.
// @Tags usuarios
// @Accept json
// @Produce json
// @Param data body map[string]string true "Datos del usuario"
// @Success 200 {object} models.Client
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users/get-or-create [post]
func GetOrCreateClient(c *gin.Context) {
	var input struct {
		Phone string `json:"phone" binding:"required"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "details": err.Error()})
		return
	}

	if input.Phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Número de teléfono requerido"})
		return
	}

	user, err := clientRepo.GetOrCreateClient(input.Phone, input.Name, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar usuario", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Usuario procesado exitosamente",
		"user":    user,
	})
}

// GetCurrentUser obtiene el usuario autenticado desde el token
// @Summary Obtener usuario actual
// @Description Retorna el perfil del usuario autenticado
// @Tags usuarios
// @Produce json
// @Success 200 {object} models.Client
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/users/me [get]
func GetCurrentUser(c *gin.Context) {
	// Obtener ID del usuario desde el middleware de autenticación
	userID, exists := c.Get("current_user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario no autenticado"})
		return
	}

	// Convertir interface{} a uint
	id, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error de autenticación"})
		return
	}

	user, err := clientRepo.GetClientByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}

	c.JSON(http.StatusOK, user)
}
