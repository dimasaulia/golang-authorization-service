package entities

import "time"

type TeamMember struct {
	ID         int64      `db:"id" json:"id"`
	TeamId     int64      `db:"team_id" json:"team_id"`
	UserId     int64      `db:"user_id" json:"user_id"`
	RoleInTeam *string    `db:"role_in_team" json:"role_in_team,omitempty"`
	CreatedAt  *time.Time `db:"created_at" json:"created_at,omitempty"`
}
