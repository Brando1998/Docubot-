package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/brando1998/docubot-api/models"
	"github.com/brando1998/docubot-api/repositories"
)

var (
	conversationRepo repositories.ConversationRepository
	botRepo          repositories.BotRepository
	clientRepo       repositories.ClientRepository
)

type IncomingMessageRequest struct {
	Phone     string `json:"phone"`
	Message   string `json:"message"`
	BotNumber string `json:"botNumber"`
}

type RasaResponseItem struct {
	Text string `json:"text"`
}

// Setters para inyección de dependencias
func SetConversationRepo(repo repositories.ConversationRepository) {
	conversationRepo = repo
}

func SetBotRepo(repo repositories.BotRepository) {
	botRepo = repo
}

func SetClientRepo(repo repositories.ClientRepository) {
	clientRepo = repo
}

// HandleWebSocket maneja conexiones WebSocket entrantes
func HandleWebSocket(c *gin.Context, hub *WebSocketHub, upgrader websocket.Upgrader) {

	//Obtener el numero
	log.Println("Nueva conexión WebSocket intentada")
	botPhone := c.Query("phone") // Este es el número DEL BOT
	if botPhone == "" {
		log.Println("Bot phone number missing in WebSocket connection")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Verificar si el bot ya está registrado
	if _, err := hub.GetBotConnection(botPhone); err == nil {
		log.Printf("Bot %s ya está registrado", botPhone)
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	// Registrar la conexión del BOT
	log.Printf("Registrando bot: %s", botPhone)
	hub.RegisterBot(botPhone, conn)

	// Manejar mensajes entrantes
	go func() {
		defer func() {
			hub.UnregisterBot(botPhone)
			conn.Close()
			log.Printf("Conexión cerrada para bot: %s", botPhone)
		}()

		for {
			// Procesar mensajes, primero convertir de json a struct
			var msg IncomingMessageRequest
			if err := conn.ReadJSON(&msg); err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Printf("Error reading message: %v", err)
				}
				break
			}
			// Procesar mensaje
			if err := processIncomingMessage(msg, hub); err != nil {
				log.Printf("Error processing message: %v", err)
			}
		}
	}()
}

// processIncomingMessage procesa los mensajes entrantes
func processIncomingMessage(msg IncomingMessageRequest, hub *WebSocketHub) error {
	// Asegurar formato consistente del botNumber
	log.Printf("Procesando mensaje de %s a bot %s: %s", msg.Phone, msg.BotNumber, msg.Message)

	// 1. Procesar cliente (guardar en DB)
	cleanPhone := strings.Split(msg.Phone, "@")[0]
	client, err := clientRepo.GetOrCreateClient(cleanPhone, "", "")
	if err != nil {
		return fmt.Errorf("failed to get/create client: %w", err)
	}

	// 2. Procesar bot (guardar en DB)
	cleanBotNumber := strings.Split(msg.BotNumber, "@")[0]
	bot, err := botRepo.GetOrCreateBot(cleanBotNumber, "Default Bot")
	if err != nil {
		return fmt.Errorf("failed to get/create bot: %w", err)
	}

	// 3. Guardar mensaje del usuario
	clientMsg := models.Message{
		ClientID:  client.ID,
		BotID:     bot.ID,
		Sender:    msg.Phone,
		Text:      msg.Message,
		Timestamp: time.Now(),
	}

	if err := conversationRepo.SaveMessage(context.TODO(), client.ID, bot.ID, clientMsg); err != nil {
		return fmt.Errorf("failed to save client message: %w", err)
	}

	// 4. Procesar con Rasa
	rasaResponses, err := sendToRasa(msg.Phone, msg.Message)
	if err != nil {
		return fmt.Errorf("rasa processing failed: %w", err)
	}
	log.Printf("Respuestas de Rasa recibidas: %+v", rasaResponses)

	// 5. Procesar respuestas
	for _, response := range rasaResponses {
		if response.Text == "" {
			continue
		}

		// Guardar respuesta del bot
		botMsg := models.Message{
			ClientID:  client.ID,
			BotID:     bot.ID,
			Sender:    "bot",
			Text:      response.Text,
			Timestamp: time.Now(),
		}
		if err := conversationRepo.SaveMessage(context.TODO(), client.ID, bot.ID, botMsg); err != nil {
			log.Printf("Failed to save bot message: %v", err)
		}

		// Enviar respuesta al cliente
		log.Printf("Enviando respuesta a bot %s para cliente %s: %s",
			msg.BotNumber, msg.Phone, response.Text)

		if err := hub.SendToBot(msg.BotNumber, map[string]interface{}{
			"to":      msg.Phone, // Cliente destino
			"message": response.Text,
		}); err != nil {
			log.Printf("Failed to send message to bot: %v", err)
		}
	}

	return nil
}

// sendToRasa envía mensajes al servidor Rasa
func sendToRasa(sender, message string) ([]RasaResponseItem, error) {
	payload := map[string]interface{}{
		"sender":  sender,
		"message": message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// req, err := http.NewRequest("POST", "http://localhost:5005/webhooks/rest/webhook", bytes.NewBuffer(body))
	req, err := http.NewRequest("POST", "http://rasa:5005/webhooks/rest/webhook", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	var responses []RasaResponseItem
	if err := json.NewDecoder(resp.Body).Decode(&responses); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return responses, nil
}
