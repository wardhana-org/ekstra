package models

import "time"

type AuthToken struct {
	ID        int64      `json:"id" db:"id"`
	SessionID int64      `json:"session_id" db:"session_id"`
	TokenHash string     `json:"-" db:"token_hash"`          // Hash of the raw token; the raw access or refresh token is never stored.
	TokenType string     `json:"token_type" db:"token_type"` // Token purpose, such as "access" or "refresh".
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at" db:"revoked_at"` // Set when the token is invalidated before its normal expiration time.
}
