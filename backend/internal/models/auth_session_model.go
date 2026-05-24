package models

import "time"

type AuthSession struct {
	ID                int64      `json:"id" db:"id"`
	UserID            int64      `json:"user_id" db:"user_id"`
	ClientType        string     `json:"client_type" db:"client_type"` // Identifies where the session came from, such as "web", "mobile", or "mobile_webview".
	DeviceName        *string    `json:"device_name" db:"device_name"` // Optional user-friendly label, such as "Chrome on Windows" or "name's iPhone".
	UserAgent         *string    `json:"user_agent" db:"user_agent"`   // Raw HTTP User-Agent header sent by the browser, app, or WebView.
	IPAddress         *string    `json:"ip_address" db:"ip_address"`   // Client IP seen by the API during login; it may be affected by proxies.
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	LastSeenAt        *time.Time `json:"last_seen_at" db:"last_seen_at"` // Updated when the session is successfully used after login.
	ExpiresAt         time.Time  `json:"expires_at" db:"expires_at"`
	AbsoluteExpiresAt time.Time  `json:"absolute_expires_at" db:"absolute_expires_at"` // Hard session expiry that refresh rotation cannot extend past.
	RevokedAt         *time.Time `json:"revoked_at" db:"revoked_at"`                   // Set when the session is intentionally invalidated, such as logout or admin revocation.
	RevokedReason     *string    `json:"revoked_reason" db:"revoked_reason"`           // Short machine-readable reason, such as "logout" or "refresh_token_reuse".
	ReuseDetectedAt   *time.Time `json:"reuse_detected_at" db:"reuse_detected_at"`     // Set when refresh token reuse is detected for this session.
}
