package entities

import "time"

type UserPermissionOverride struct {
	ID           int64      `db:"id" json:"id"`
	UserId       int64      `db:"user_id" json:"user_id"`
	PermissionId int64      `db:"permission_id" json:"permission_id"`
	Effect       string     `db:"effect" json:"effect"`
	Reason       string     `db:"reason" json:"reason"`
	ExpiresAt    *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	CreatedBy    *int64     `db:"created_by" json:"created_by,omitempty"`
	CreatedAt    *time.Time `db:"created_at" json:"created_at,omitempty"`
}
