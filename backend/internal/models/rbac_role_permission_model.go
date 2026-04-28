package models

import "time"

type RBACRolePermission struct {
	ID           int64     `json:"id" db:"id"`
	RoleID       int64     `json:"role_id" db:"role_id"`
	PermissionID int64     `json:"permission_id" db:"permission_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
