package dto

type CreateTeamRoleRequest struct {
	TeamId int64  `json:"team_id"`
	RoleId int64  `json:"role_id"`
	AppId  *int64 `json:"app_id"`
}

type UpdateTeamRoleRequest struct {
	TeamId *int64 `json:"team_id"`
	RoleId *int64 `json:"role_id"`
	AppId  *int64 `json:"app_id"`
}
