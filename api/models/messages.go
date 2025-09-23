package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ClientID  uint               `bson:"client_id"`
	BotID     uint               `bson:"bot_id"`
	Sender    string             `bson:"sender"` // "user" o "bot"
	Text      string             `bson:"text"`
	Timestamp time.Time          `bson:"timestamp"`
	SessionID string             `bson:"session_id,omitempty"`
}

type Conversation struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    uint               `bson:"client_id"`
	BotID     uint               `bson:"bot_id"`
	Messages  []Message          `bson:"messages"`
	CreatedAt time.Time          `bson:"created_at"`
}

type Document struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ClientID  uint               `bson:"client_id"`
	FileName  string             `bson:"file_name"`
	URL       string             `bson:"url"`
	Type      string             `bson:"type"` // Ej: "manifiesto", "certificado"
	CreatedAt time.Time          `bson:"created_at"`
}
