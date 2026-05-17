package dto

type CreateTeamMemberRequest struct {
	TeamId     int64   `json:"team_id"`
	UserId     int64   `json:"user_id"`
	RoleInTeam *string `json:"role_in_team"`
}

type UpdateTeamMemberRequest struct {
	TeamId     *int64  `json:"team_id"`
	UserId     *int64  `json:"user_id"`
	RoleInTeam *string `json:"role_in_team"`
}
