package models

import "time"

type AuthToken struct {
	ID                  int64      `json:"id" db:"id"`
	SessionID           int64      `json:"session_id" db:"session_id"`
	TokenHash           string     `json:"-" db:"token_hash"`          // Hash of the raw token; the raw access or refresh token is never stored.
	TokenType           string     `json:"token_type" db:"token_type"` // Token purpose, such as "access" or "refresh".
	Generation          *int       `json:"generation" db:"generation"` // Refresh-token chain number, such as R1 = 1 and R2 = 2; access tokens leave this empty.
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt           time.Time  `json:"expires_at" db:"expires_at"`
	UsedAt              *time.Time `json:"used_at" db:"used_at"`                               // Set when a refresh token is consumed to issue the next refresh token.
	ReplacedByTokenID   *int64     `json:"replaced_by_token_id" db:"replaced_by_token_id"`     // Points from the old refresh token to the new refresh token that replaced it.
	ReuseGraceExpiresAt *time.Time `json:"reuse_grace_expires_at" db:"reuse_grace_expires_at"` // Short retry window for the immediately previous refresh token.
	RevokedAt           *time.Time `json:"revoked_at" db:"revoked_at"`                         // Set when the token is invalidated before its normal expiration time.
	RevokedReason       *string    `json:"revoked_reason" db:"revoked_reason"`                 // Short reason, such as "rotated", "logout", or "refresh_token_reuse".
}
