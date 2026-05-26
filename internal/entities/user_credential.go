package entities

import "time"

type UserCredential struct {
	ID                 int64      `db:"id" json:"id"`
	UserID             int64      `db:"user_id" json:"user_id"`
	PasswordHash       string     `db:"password_hash" json:"-"`
	MustChangePassword bool       `db:"must_change_password" json:"must_change_password"`
	PasswordChangedAt  *time.Time `db:"password_changed_at" json:"password_changed_at,omitempty"`
	CreatedAt          *time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt          *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
