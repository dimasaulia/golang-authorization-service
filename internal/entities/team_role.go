package entities

import "time"

type TeamRole struct {
	ID        int64      `db:"id" json:"id"`
	TeamId    int64      `db:"team_id" json:"team_id"`
	RoleId    int64      `db:"role_id" json:"role_id"`
	AppId     *int64     `db:"app_id" json:"app_id,omitempty"`
	CreatedAt *time.Time `db:"created_at" json:"created_at,omitempty"`
}
