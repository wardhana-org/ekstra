package models

import "time"

type AuthSecurityEvent struct {
	ID        int64     `json:"id" db:"id"`
	UserID    *int64    `json:"user_id" db:"user_id"`
	SessionID *int64    `json:"session_id" db:"session_id"`
	TokenID   *int64    `json:"token_id" db:"token_id"`
	EventType string    `json:"event_type" db:"event_type"`
	IPAddress *string   `json:"ip_address" db:"ip_address"`
	UserAgent *string   `json:"user_agent" db:"user_agent"`
	Metadata  []byte    `json:"metadata" db:"metadata"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
