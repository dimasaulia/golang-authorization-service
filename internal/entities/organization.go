package entities

import "time"

type Organization struct {
	ID        int64      `db:"id" json:"id"`
	Code      string     `db:"code" json:"code"`
	Name      string     `db:"name" json:"name"`
	Type      string     `db:"type" json:"type"`
	Status    string     `db:"status" json:"status"`
	CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
