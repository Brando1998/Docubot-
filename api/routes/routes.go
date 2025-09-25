package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/brando1998/docubot-api/controllers"
	_ "github.com/brando1998/docubot-api/docs"
	"github.com/brando1998/docubot-api/middleware"
)

type RouterConfig struct {
	WSHub    *controllers.WebSocketHub
	Upgrader *websocket.Upgrader
}

func SetupRoutes(r *gin.Engine, config *RouterConfig) {
	// =============================================
	// Middlewares globales
	// =============================================
	r.Use(
		middleware.LoggerMiddleware(),
		middleware.CORSMiddleware(),
	)

	// =============================================
	// Rutas Públicas (sin autenticación)
	// =============================================
	public := r.Group("/")
	{
		public.GET("/health", controllers.Health)
		public.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// WebSocket para conexión con Baileys
		public.GET("/ws", func(c *gin.Context) {
			controllers.HandleWebSocket(c, config.WSHub, *config.Upgrader)
		})

		// Debug: listar bots conectados
		public.GET("/debug/bots", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"bots":  config.WSHub.ListBots(),
				"total": len(config.WSHub.ListBots()),
			})
		})
	}

	// =============================================
	// Autenticación
	// =============================================
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", controllers.LoginWithPaseto)
		authGroup.POST("/refresh", controllers.RefreshPasetoToken)
	}

	// =============================================
	// Rutas Protegidas con PASETO
	// =============================================
	api := r.Group("/api/v1")
	api.Use(middleware.PasetoAuthMiddleware())
	{
		// --------------------------
		// Usuarios
		// --------------------------
		userGroup := api.Group("/users")
		{
			userGroup.GET("/me", controllers.GetCurrentUser)                // Usuario actual
			userGroup.POST("", controllers.CreateClient)                    // Crear usuario
			userGroup.GET("/id/:id", controllers.GetClientByID)             // Obtener por ID
			userGroup.GET("/phone/:phone", controllers.GetClientByPhone)    // Obtener por teléfono
			userGroup.POST("/get-or-create", controllers.GetOrCreateClient) // Obtener o crear
		}

		//--------------------------
		//Conversaciones (ya implementado)
		//--------------------------
		// convGroup := api.Group("/conversations")
		// {
		// 	// Nota: Los mensajes se manejan via WebSocket
		// 	// Aquí se pueden agregar endpoints para historial
		// 	convGroup.GET("/history/:client_id", controllers.GetConversationHistory)
		// 	convGroup.GET("/recent", controllers.GetRecentConversations)
		// }

		// --------------------------
		// WhatsApp (administración y dashboard)
		// --------------------------
		whatsappGroup := api.Group("/whatsapp")
		{
			// Endpoints para dashboard (sin restricción admin)
			whatsappGroup.GET("/qr", controllers.GetWhatsAppQR)
			whatsappGroup.GET("/status", controllers.GetSessionStatus)

			// Endpoints para administradores
			whatsappGroup.POST("/send", controllers.SendWhatsAppMessage)
			whatsappGroup.GET("/session/:session_id", controllers.GetWhatsAppSession)
			whatsappGroup.POST("/session", controllers.CreateWhatsAppSession)
		}
	}

	// =============================================
	// Rutas de Administración (futuro)
	// =============================================
	// admin := r.Group("/admin")
	// admin.Use(middleware.PasetoAdminMiddleware()) // Middleware para admin
	// {
	// 	admin.GET("/metrics", controllers.GetMetrics)
	// 	admin.GET("/users", controllers.GetAllUsers)
	// 	admin.DELETE("/users/:id", controllers.DeleteUser)
	// }
}
