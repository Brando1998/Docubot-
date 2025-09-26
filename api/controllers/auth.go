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

// ❌ PROBLEMA: Esta línea causaba el nil pointer porque se ejecuta antes de ConnectPostgres()
// var authDB = database.GetDB()

// PasetoPayload define la estructura del payload del token PASETO
type PasetoPayload struct {
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

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

	// ✅ SOLUCIÓN: Obtener DB directamente en la función
	db := database.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error de conexión a base de datos"})
		return
	}

	// Buscar usuario en la base de datos
	var user models.SystemUser
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	// Verificar si el usuario está activo
	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario inactivo"})
		return
	}

	// Verificar contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credenciales inválidas"})
		return
	}

	// Actualizar último login
	now := time.Now()
	user.LastLogin = &now
	db.Save(&user)

	// Generar token PASETO
	token, expiresAt, err := generatePasetoToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error generando token"})
		return
	}

	// Respuesta de login exitoso
	response := LoginResponse{
		AccessToken: token,
		ExpiresAt:   expiresAt,
	}
	response.User.ID = user.ID
	response.User.Username = user.Username
	response.User.Email = user.Email
	response.User.Role = user.Role

	c.JSON(http.StatusOK, response)
}

// generatePasetoToken genera un token PASETO
func generatePasetoToken(userID uint, username, role string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(tokenDuration)

	payload := PasetoPayload{
		UserID:    userID,
		Username:  username,
		Role:      role,
		IssuedAt:  now,
		ExpiresAt: expiresAt,
	}

	v2 := paseto.NewV2()
	secretKey := []byte(os.Getenv("PASETO_SECRET_KEY"))

	// Verificar que la clave secreta esté configurada
	if len(secretKey) == 0 {
		secretKey = []byte("default-secret-key-change-in-production-32-chars")
	}

	token, err := v2.Encrypt(secretKey, payload, nil)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expiresAt, nil
}

// RefreshPasetoToken maneja la renovación de tokens
func RefreshPasetoToken(c *gin.Context) {
	// Obtener token del header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token no proporcionado"})
		return
	}

	// Extraer token
	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	// Verificar token actual
	payload, err := verifyPasetoToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token inválido"})
		return
	}

	// ✅ SOLUCIÓN: Obtener DB directamente en la función
	db := database.GetDB()
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error de conexión a base de datos"})
		return
	}

	// Verificar que el usuario aún existe
	var user models.SystemUser
	if err := db.First(&user, payload.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no encontrado"})
		return
	}

	// Generar nuevo token
	newToken, expiresAt, err := generatePasetoToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error generando nuevo token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newToken,
		"expires_at":   expiresAt,
	})
}

// verifyPasetoToken verifica y decodifica el token
func verifyPasetoToken(token string) (*PasetoPayload, error) {
	v2 := paseto.NewV2()
	var payload PasetoPayload

	secretKey := []byte(os.Getenv("PASETO_SECRET_KEY"))
	if len(secretKey) == 0 {
		secretKey = []byte("default-secret-key-change-in-production-32-chars")
	}

	if err := v2.Decrypt(token, secretKey, &payload, nil); err != nil {
		return nil, errors.New("token inválido")
	}

	if time.Now().After(payload.ExpiresAt) {
		return nil, errors.New("token expirado")
	}

	return &payload, nil
}
