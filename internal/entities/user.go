package entities

import "time"

type User struct {
	ID              int64      `db:"id" json:"id"`
	OrganizationId  int64      `db:"organization_id" json:"organization_id"`
	Username        string     `db:"username" json:"username"`
	Email           string     `db:"email" json:"email"`
	DisplayName     string     `db:"display_name" json:"display_name"`
	Type            string     `db:"type" json:"type"`
	Status          string     `db:"status" json:"status"`
	EmailVerifiedAt *time.Time `db:"email_verified_at" json:"email_verified_at,omitempty"`
	CreatedAt       *time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
