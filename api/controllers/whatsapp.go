package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// WhatsAppQRResponse estructura para la respuesta del QR
type WhatsAppQRResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	QRCode    string `json:"qr_code,omitempty"`
	QRImage   string `json:"qr_image,omitempty"`
	Connected bool   `json:"connected"`
}

// WhatsAppStatusResponse estructura para el estado de la sesión
type WhatsAppStatusResponse struct {
	Status      string    `json:"status"`
	Message     string    `json:"message,omitempty"`
	Connected   bool      `json:"connected"`
	BotNumber   string    `json:"bot_number,omitempty"`
	LastSeen    time.Time `json:"last_seen,omitempty"`
	SessionInfo struct {
		Name   string `json:"name,omitempty"`
		Avatar string `json:"avatar,omitempty"`
	} `json:"session_info"`
}

// BaileysStatusRequest estructura para consultar estado a Baileys
type BaileysStatusRequest struct {
	Action string `json:"action"`
}

// GetWhatsAppQR obtiene el código QR de WhatsApp desde Baileys
// @Summary Obtener código QR de WhatsApp
// @Description Retorna el código QR para conectar WhatsApp o el estado de conexión actual
// @Tags whatsapp
// @Produce json
// @Success 200 {object} WhatsAppQRResponse
// @Failure 500 {object} map[string]string
// @Router /whatsapp/qr [get]
func GetWhatsAppQR(c *gin.Context) {
	// URL del servicio Baileys
	baileysURL := "http://baileys:3000/qr" // Ajusta según tu configuración

	// Hacer solicitud a Baileys para obtener QR
	resp, err := http.Get(baileysURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "No se pudo conectar con el servicio de WhatsApp",
			"details": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var qrResponse WhatsAppQRResponse
	if err := json.NewDecoder(resp.Body).Decode(&qrResponse); err != nil {
		// Si Baileys no responde con JSON esperado, verificar estado manualmente
		status := checkBaileysConnection()
		if status.Connected {
			c.JSON(http.StatusOK, WhatsAppQRResponse{
				Status:    "connected",
				Message:   "WhatsApp ya está conectado",
				Connected: true,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al procesar respuesta del servicio WhatsApp",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, qrResponse)
}

// GetSessionStatus obtiene el estado actual de la sesión de WhatsApp
// @Summary Obtener estado de sesión WhatsApp
// @Description Retorna información sobre el estado de la conexión de WhatsApp
// @Tags whatsapp
// @Produce json
// @Success 200 {object} WhatsAppStatusResponse
// @Failure 500 {object} map[string]string
// @Router /whatsapp/status [get]
func GetSessionStatus(c *gin.Context) {
	status := checkBaileysConnection()
	c.JSON(http.StatusOK, status)
}

// checkBaileysConnection verifica el estado de conexión con Baileys
func checkBaileysConnection() WhatsAppStatusResponse {
	// Intentar conectar con el endpoint de estado de Baileys
	baileysStatusURL := "http://baileys:3000/status"

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(baileysStatusURL)

	if err != nil {
		return WhatsAppStatusResponse{
			Status:    "disconnected",
			Connected: false,
			Message:   "No se pudo conectar con el servicio Baileys",
		}
	}
	defer resp.Body.Close()

	var statusResponse struct {
		Status  string `json:"status"`
		Uptime  int    `json:"uptime"`
		BotID   string `json:"bot_id,omitempty"`
		BotName string `json:"bot_name,omitempty"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&statusResponse); err != nil {
		return WhatsAppStatusResponse{
			Status:    "unknown",
			Connected: false,
			Message:   "Error al leer estado del servicio",
		}
	}

	// Determinar si está conectado basado en la respuesta
	connected := statusResponse.Status == "running" && statusResponse.BotID != ""

	response := WhatsAppStatusResponse{
		Status:    statusResponse.Status,
		Connected: connected,
		BotNumber: statusResponse.BotID,
		LastSeen:  time.Now(),
	}

	if connected {
		response.SessionInfo.Name = statusResponse.BotName
	}

	return response
}

// SendWhatsAppMessage envía un mensaje a través de WhatsApp (para administradores)
// @Summary Enviar mensaje por WhatsApp
// @Description Envía un mensaje a un número específico a través de WhatsApp
// @Tags whatsapp
// @Accept json
// @Produce json
// @Param message body map[string]string true "Datos del mensaje"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /whatsapp/send [post]
func SendWhatsAppMessage(c *gin.Context) {
	var request struct {
		To      string `json:"to" binding:"required"`
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos"})
		return
	}

	// Preparar mensaje para Baileys
	messagePayload := map[string]string{
		"to":      request.To,
		"message": request.Message,
	}

	jsonData, err := json.Marshal(messagePayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al preparar mensaje"})
		return
	}

	// Enviar a Baileys
	baileysURL := "http://baileys:3000/send"
	resp, err := http.Post(baileysURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error al enviar mensaje",
			"details": err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en el servicio de WhatsApp"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Mensaje enviado correctamente",
		"to":      request.To,
	})
}

// CreateWhatsAppSession crea una nueva sesión de WhatsApp (placeholder para futuras funcionalidades)
func CreateWhatsAppSession(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error":   "Funcionalidad no implementada aún",
		"message": "La creación de sesiones se maneja automáticamente por Baileys",
	})
}

// GetWhatsAppSession obtiene información de una sesión específica (placeholder)
func GetWhatsAppSession(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de sesión requerido"})
		return
	}

	// Por ahora retornamos el estado general
	status := checkBaileysConnection()
	c.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
		"status":     status,
	})
}
