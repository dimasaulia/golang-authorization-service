package entities

import "time"

type Team struct {
	ID             int64      `db:"id" json:"id"`
	OrganizationId int64      `db:"organization_id" json:"organization_id"`
	Code           string     `db:"code" json:"code"`
	Name           string     `db:"name" json:"name"`
	Status         string     `db:"status" json:"status"`
	CreatedAt      *time.Time `db:"created_at" json:"created_at,omitempty"`
}
