package models

import "time"

type PlatformOperator struct {
	ID         int64     `json:"id" db:"id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	PlatformID int64     `json:"platform_id" db:"platform_id"`
	RoleID     int64     `json:"role_id" db:"role_id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
