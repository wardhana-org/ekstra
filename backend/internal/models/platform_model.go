package models

import "time"

type Platform struct {
	ID          int64     `json:"id" db:"id"`
	OwnerUserID int64     `json:"owner_user_id" db:"owner_user_id"`
	Name        string    `json:"name" db:"name"`
	Slug        string    `json:"slug" db:"slug"`
	Description *string   `json:"description" db:"description"`
	IsPublic    bool      `json:"is_public" db:"is_public"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
