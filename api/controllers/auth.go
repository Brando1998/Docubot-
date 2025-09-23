package controllers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
	"github.com/o1egl/paseto"

	database "github.com/brando1998/docubot-api/databases"
	"github.com/brando1998/docubot-api/models"
)

const (
	tokenDuration = 24 * time.Hour
)

var authDB = database.GetDB()

// LoginRequest define la estructura para el request de login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse define la estructura de respuesta para el login
type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	User        struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	} `json:"user"`
}

// LoginWithPaseto maneja el inicio de sesión
func LoginWithPaseto(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos inválidos"})
		return
	}

	// Buscar usuario en la base de datos
	var user models.SystemUser
	if err := authDB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	// Generar token PASETO
	token, payload, err := generatePasetoToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no se pudo generar el token"})
		return
	}

	// Actualizar último login
	now := time.Now()
	authDB.Model(&user).Update("last_login", &now)

	// Devolver respuesta
	response := LoginResponse{
		AccessToken: token,
		ExpiresAt:   payload.ExpiresAt,
	}
	response.User.ID = user.ID
	response.User.Username = user.Username
	response.User.Email = user.Email
	response.User.Role = user.Role

	c.JSON(http.StatusOK, response)
}

// RefreshPasetoToken maneja la renovación de tokens
func RefreshPasetoToken(c *gin.Context) {
	// Implementación similar a LoginWithPaseto pero verificando un refresh token
}

// generatePasetoToken genera un nuevo token PASETO
func generatePasetoToken(userID uint, role string) (string, *PasetoPayload, error) {
	v2 := paseto.NewV2()

	payload := &PasetoPayload{
		UserID:    userID,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(tokenDuration),
	}

	secretKey := []byte(os.Getenv("PASETO_SECRET_KEY"))
	if len(secretKey) == 0 {
		return "", nil, errors.New("PASETO_SECRET_KEY no configurada")
	}

	token, err := v2.Encrypt(secretKey, payload, nil)
	if err != nil {
		return "", nil, err
	}

	return token, payload, nil
}

// PasetoPayload define la estructura del payload del token
type PasetoPayload struct {
	UserID    uint      `json:"user_id"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
