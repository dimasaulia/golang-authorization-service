package entities

import "time"

type UserRole struct {
	ID             int64      `db:"id" json:"id"`
	UserId         int64      `db:"user_id" json:"user_id"`
	RoleId         int64      `db:"role_id" json:"role_id"`
	AppId          *int64     `db:"app_id" json:"app_id,omitempty"`
	OrganizationId *int64     `db:"organization_id" json:"organization_id,omitempty"`
	ExpiresAt      *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	AssignedBy     *int64     `db:"assigned_by" json:"assigned_by,omitempty"`
	CreatedAt      *time.Time `db:"created_at" json:"created_at,omitempty"`
}
