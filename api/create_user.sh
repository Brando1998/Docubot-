#!/bin/bash

# Docubot - Script de Creación de Usuario
# =====================================

set -e

# Colores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🤖 Docubot - Creador de Usuario${NC}"
echo -e "${BLUE}================================${NC}"
echo ""

# Verificar que estamos en el contenedor correcto
if [ ! -f "/app/go.mod" ]; then
    echo -e "${RED}❌ Error: Este script debe ejecutarse dentro del contenedor API${NC}"
    echo -e "${YELLOW}💡 Uso: make dev-shell-api${NC}"
    exit 1
fi

# Crear el archivo Go si no existe
if [ ! -f "/app/cmd/create_user.go" ]; then
    echo -e "${YELLOW}📝 Creando archivo create_user.go...${NC}"
    
    mkdir -p /app/cmd
    
    cat > /app/cmd/create_user.go << 'EOF'
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"

	database "github.com/brando1998/docubot-api/databases"
	"github.com/brando1998/docubot-api/models"
)

func main() {
	fmt.Println("🤖 Docubot - Registro de Usuario del Sistema")
	fmt.Println("===========================================")

	// Inicializar conexión a la base de datos
	database.Connect()
	db := database.GetDB()

	// Migrar el modelo si es necesario
	err := db.AutoMigrate(&models.SystemUser{})
	if err != nil {
		log.Fatalf("❌ Error al migrar la base de datos: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)

	// Solicitar datos del usuario
	fmt.Print("📝 Nombre de usuario: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	if username == "" {
		log.Fatal("❌ El nombre de usuario es obligatorio")
	}

	// Verificar si el usuario ya existe
	var existingUser models.SystemUser
	if err := db.Where("username = ?", username).First(&existingUser).Error; err == nil {
		log.Fatal("❌ El usuario ya existe")
	}

	fmt.Print("📧 Email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	if email == "" {
		log.Fatal("❌ El email es obligatorio")
	}

	// Verificar si el email ya existe
	if err := db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		log.Fatal("❌ El email ya existe")
	}

	// Solicitar contraseña de forma segura
	fmt.Print("🔒 Contraseña: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatal("❌ Error al leer la contraseña")
	}
	password := string(passwordBytes)
	fmt.Println() // Nueva línea después de la contraseña

	if len(password) < 6 {
		log.Fatal("❌ La contraseña debe tener al menos 6 caracteres")
	}

	fmt.Print("👤 Rol (admin/user) [user]: ")
	role, _ := reader.ReadString('\n')
	role = strings.TrimSpace(role)

	if role == "" {
		role = "user"
	}

	if role != "admin" && role != "user" {
		log.Fatal("❌ El rol debe ser 'admin' o 'user'")
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("❌ Error al hash la contraseña: %v", err)
	}

	// Crear usuario
	user := models.SystemUser{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		IsActive:     true,
	}

	// Guardar en la base de datos
	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("❌ Error al crear el usuario: %v", err)
	}

	fmt.Printf("\n✅ Usuario creado exitosamente!\n")
	fmt.Printf("📝 Username: %s\n", user.Username)
	fmt.Printf("📧 Email: %s\n", user.Email)
	fmt.Printf("👤 Rol: %s\n", user.Role)
	fmt.Printf("🆔 ID: %d\n", user.ID)
	fmt.Printf("📅 Creado: %s\n", user.CreatedAt.Format("2006-01-02 15:04:05"))

	fmt.Println("\n🎉 ¡Ahora puedes hacer login en el dashboard!")
}
EOF

    echo -e "${GREEN}✅ Archivo create_user.go creado${NC}"
fi

# Compilar y ejecutar
echo -e "${YELLOW}🔨 Compilando comando...${NC}"
cd /app

# Compilar el comando
go build -o /tmp/create_user ./cmd/create_user.go

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Error al compilar${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Compilado exitosamente${NC}"
echo ""

# Ejecutar el comando
echo -e "${BLUE}🚀 Ejecutando creador de usuario...${NC}"
echo ""
/tmp/create_user