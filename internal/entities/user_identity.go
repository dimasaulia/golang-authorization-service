package entities

import "time"

type UserIdentity struct {
	ID             int64      `db:"id" json:"id"`
	UserId         int64      `db:"user_id" json:"user_id"`
	Provider       string     `db:"provider" json:"provider"`
	ProviderUserId string     `db:"provider_user_id" json:"provider_user_id"`
	Username       *string    `db:"username" json:"username,omitempty"`
	Email          *string    `db:"email" json:"email,omitempty"`
	IsPrimary      bool       `db:"is_primary" json:"is_primary"`
	CreatedAt      *time.Time `db:"created_at" json:"created_at,omitempty"`
}
