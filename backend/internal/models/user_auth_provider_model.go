package models

import "time"

type UserAuthProvider struct {
	ID             int64     `json:"id" db:"id"`
	UserID         int64     `json:"user_id" db:"user_id"`
	Provider       string    `json:"provider" db:"provider"`
	ProviderUserID *string   `json:"provider_user_id" db:"provider_user_id"`
	PasswordHash   *string   `json:"-" db:"password_hash"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
