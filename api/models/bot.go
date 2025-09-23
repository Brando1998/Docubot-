package models

type Bot struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	Type   string `json:"type"`   // transporte, salud, etc.
	Number string `json:"number"` // Número de WhatsApp asociado
	Active bool   `json:"active"`
}
