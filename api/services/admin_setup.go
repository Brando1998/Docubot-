package services

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/brando1998/docubot-api/models"
)

// DefaultAdminCredentials contiene las credenciales por defecto del admin
type DefaultAdminCredentials struct {
	Username string
	Email    string
	Password string
}

// GetDefaultAdminCredentials obtiene las credenciales del admin desde variables de entorno o usa valores por defecto
func GetDefaultAdminCredentials() DefaultAdminCredentials {
	return DefaultAdminCredentials{
		Username: getEnvOrDefault("ADMIN_USERNAME", "admin"),
		Email:    getEnvOrDefault("ADMIN_EMAIL", "admin@docubot.local"),
		Password: getEnvOrDefault("ADMIN_PASSWORD", "DocubotAdmin123!"),
	}
}

// EnsureDefaultAdminUser verifica si existe un usuario administrador y lo crea si no existe
func EnsureDefaultAdminUser(db *gorm.DB) error {
	log.Println("🔍 Verificando usuario administrador por defecto...")

	// Verificar si ya existe un usuario con rol admin
	var adminExists int64
	err := db.Model(&models.SystemUser{}).Where("role = ?", "admin").Count(&adminExists).Error
	if err != nil {
		return err
	}

	// Si ya existe un admin, no hacer nada
	if adminExists > 0 {
		log.Printf("✅ Usuario administrador ya existe (%d admins encontrados)", adminExists)
		return nil
	}

	// Obtener credenciales por defecto
	creds := GetDefaultAdminCredentials()

	// Verificar si ya existe un usuario con el mismo username o email
	var existingUser models.SystemUser
	err = db.Where("username = ? OR email = ?", creds.Username, creds.Email).First(&existingUser).Error
	if err == nil {
		log.Printf("⚠️  Usuario con username '%s' o email '%s' ya existe pero no es admin. Saltando creación automática.",
			creds.Username, creds.Email)
		return nil
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Crear usuario administrador
	adminUser := models.SystemUser{
		Username:     creds.Username,
		Email:        creds.Email,
		PasswordHash: string(hashedPassword),
		Role:         "admin",
		IsActive:     true,
	}

	// Guardar en la base de datos
	if err := db.Create(&adminUser).Error; err != nil {
		return err
	}

	log.Printf("🎉 ¡Usuario administrador creado exitosamente!")
	log.Printf("📝 Username: %s", adminUser.Username)
	log.Printf("📧 Email: %s", adminUser.Email)
	log.Printf("👤 Rol: %s", adminUser.Role)
	log.Printf("🆔 ID: %d", adminUser.ID)

	// Log de credenciales para development (solo si no son variables de entorno personalizadas)
	if os.Getenv("ADMIN_USERNAME") == "" && os.Getenv("ADMIN_PASSWORD") == "" {
		log.Printf("🔑 CREDENCIALES POR DEFECTO:")
		log.Printf("   Username: %s", creds.Username)
		log.Printf("   Password: %s", creds.Password)
		log.Printf("⚠️  CAMBIA ESTAS CREDENCIALES EN PRODUCCIÓN usando variables de entorno!")
	}

	return nil
}

// getEnvOrDefault obtiene una variable de entorno o retorna un valor por defecto
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
