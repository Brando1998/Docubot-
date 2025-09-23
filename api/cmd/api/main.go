// @title Docubot API
// @version 1.0
// @description API para generación automatizada de documentos y chatbot.
// @host localhost:8080
// @BasePath /
// @schemes http

package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/brando1998/docubot-api/config"
	"github.com/brando1998/docubot-api/controllers"
	database "github.com/brando1998/docubot-api/databases"
	"github.com/brando1998/docubot-api/models"
	"github.com/brando1998/docubot-api/repositories"
	"github.com/brando1998/docubot-api/routes"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func initDependencies() (*controllers.WebSocketHub, *gin.Engine) {
	// 1. Configuración inicial
	config.LoadEnv()

	// 2. Conexiones a bases de datos
	if err := database.ConnectPostgres(); err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	if err := database.ConnectMongoDB(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// 3. Migraciones
	runMigrations()

	// 4. Inicialización de repositorios
	initRepositories()

	// 5. WebSocket Hub
	wsHub := controllers.NewWebSocketHub()

	// 6. Configuración de Gin
	routerConfig := &routes.RouterConfig{
		WSHub:    wsHub,
		Upgrader: &upgrader,
	}
	r := gin.Default()

	routes.SetupRoutes(r, routerConfig)

	return wsHub, r
}

func runMigrations() {
	err := database.DB.AutoMigrate(
		&models.Client{},
		&models.Bot{},
		&models.WhatsAppSession{},
	)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}

func initRepositories() {
	conversationRepo := repositories.NewConversationRepository(database.MongoClient)
	clientRepo := repositories.NewClientRepository(database.DB)
	botRepo := repositories.NewBotRepository(database.DB)

	controllers.SetConversationRepo(conversationRepo)
	controllers.SetClientRepo(clientRepo)
	controllers.SetBotRepo(botRepo)
}

func getServerPort() string {
	if port := os.Getenv("PORT"); port != "" {
		return port
	}
	return "8080" // Default port
}

func main() {
	_, router := initDependencies()

	port := getServerPort()
	log.Printf("🚀 Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
