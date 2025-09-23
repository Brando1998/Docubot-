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

		// WebSocket
		public.GET("/ws", func(c *gin.Context) {
			controllers.HandleWebSocket(c, config.WSHub, *config.Upgrader)
		})
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
			userGroup.POST("", controllers.CreateClient)
			userGroup.GET("/id/:id", controllers.GetClientByID)
			userGroup.GET("/phone/:phone", controllers.GetClientByPhone)
			userGroup.POST("/get-or-create", controllers.GetOrCreateClient)
		}

		// --------------------------
		// Conversaciones
		// --------------------------
		// convGroup := api.Group("/conversations")
		// {
		// 	convGroup.POST("", controllers.SendMessage)
		// 	convGroup.GET("/history", controllers.GetConversationHistory)
		// }

		// --------------------------
		// Documentos
		// --------------------------
		docGroup := api.Group("/documents")
		{
			docGroup.POST("", controllers.UploadDocument)
			docGroup.GET("", controllers.GetUserDocuments)
			// docGroup.GET("/:doc_id", controllers.GetDocument)
		}

		// --------------------------
		// Automatización
		// --------------------------
		autoGroup := api.Group("/automation")
		{
			autoGroup.POST("", controllers.RunAutomation)
			// autoGroup.GET("/status/:job_id", controllers.GetAutomationStatus)
		}

		// --------------------------
		// WhatsApp (solo para administración)
		// --------------------------
		// whatsappGroup := api.Group("/whatsapp")
		// whatsappGroup.Use(middleware.PasetoAdminMiddleware())
		// {
		// 	whatsappGroup.POST("/send", controllers.SendWhatsAppMessage)
		// 	whatsappGroup.GET("/session/:user_id", controllers.GetWhatsAppSession)
		// 	whatsappGroup.POST("/session", controllers.CreateWhatsAppSession)
		// }
	}

	// =============================================
	// Rutas de Administración
	// =============================================
	// admin := r.Group("/admin")
	// admin.Use(middleware.PasetoAdminMiddleware())
	// {
	// 	admin.GET("/metrics", controllers.GetMetrics)
	// }
}
