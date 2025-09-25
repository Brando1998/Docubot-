package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// Solo cargar .env si no estamos en contenedor
	if os.Getenv("DOCKER_CONTAINER") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("No se pudo cargar .env, usando variables del entorno.")
		} else {
			log.Println("✅ Archivo .env cargado para desarrollo local")
		}
	} else {
		log.Println("🐳 Ejecutándose en contenedor, usando variables de entorno")
	}
}
