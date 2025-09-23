package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// WhatsAppSession represents a WhatsApp session for a user
type WhatsAppSession struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ClientID  uint               `bson:"client_id"`
	SessionID string             `bson:"session_id"`
	Status    string             `bson:"status"` // online, offline, etc.
	CreatedAt time.Time          `bson:"created_at"`
}
