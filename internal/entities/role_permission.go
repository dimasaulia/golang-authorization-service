package entities

import "time"

type RolePermission struct {
	ID           int64      `db:"id" json:"id"`
	RoleId       int64      `db:"role_id" json:"role_id"`
	PermissionId int64      `db:"permission_id" json:"permission_id"`
	Effect       string     `db:"effect" json:"effect"`
	CreatedAt    *time.Time `db:"created_at" json:"created_at,omitempty"`
}
