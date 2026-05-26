package entities

import "time"

type UserVerificationCode struct {
	ID        int64      `db:"id" json:"id"`
	UserID    int64      `db:"user_id" json:"user_id"`
	Purpose   string     `db:"purpose" json:"purpose"`
	CodeHash  string     `db:"code_hash" json:"-"`
	ExpiresAt time.Time  `db:"expires_at" json:"expires_at"`
	UsedAt    *time.Time `db:"used_at" json:"used_at,omitempty"`
	CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
}
