package models

import "time"

// PasetoPayload define la estructura del payload del token PASETO
type PasetoPayload struct {
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
