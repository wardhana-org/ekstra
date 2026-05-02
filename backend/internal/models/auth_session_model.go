package models

import "time"

type AuthSession struct {
	ID         int64      `json:"id" db:"id"`
	UserID     int64      `json:"user_id" db:"user_id"`
	ClientType string     `json:"client_type" db:"client_type"`
	DeviceName *string    `json:"device_name" db:"device_name"`
	UserAgent  *string    `json:"user_agent" db:"user_agent"`
	IPAddress  *string    `json:"ip_address" db:"ip_address"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	LastSeenAt *time.Time `json:"last_seen_at" db:"last_seen_at"`
	ExpiresAt  time.Time  `json:"expires_at" db:"expires_at"`
	RevokedAt  *time.Time `json:"revoked_at" db:"revoked_at"`
}
