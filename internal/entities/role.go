package entities

import "time"

type Role struct {
	ID             int64      `db:"id" json:"id"`
	OrganizationId *int64     `db:"organization_id" json:"organization_id,omitempty"`
	AppId          *int64     `db:"app_id" json:"app_id,omitempty"`
	Code           string     `db:"code" json:"code"`
	Name           string     `db:"name" json:"name"`
	Description    *string    `db:"description" json:"description,omitempty"`
	Scope          string     `db:"scope" json:"scope"`
	IsSystem       bool       `db:"is_system" json:"is_system"`
	Status         string     `db:"status" json:"status"`
	CreatedAt      *time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt      *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
