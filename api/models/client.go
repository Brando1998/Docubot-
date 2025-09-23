package models

import (
	"time"
)

type Client struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Email     *string   `json:"email" gorm:"uniqueIndex"` // Cambiado a puntero para permitir nil
	Phone     string    `json:"phone" gorm:"uniqueIndex"` // importante para identificar desde WhatsApp
	Company   string    `json:"company"`                  // para cuando agregues login
	BotID     uint      `json:"bot_id"`                   // para cuando agregues login
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
