package dto

type CreateRolePermissionRequest struct {
	RoleId       int64  `json:"role_id"`
	PermissionId int64  `json:"permission_id"`
	Effect       string `json:"effect"`
}

type UpdateRolePermissionRequest struct {
	RoleId       *int64  `json:"role_id"`
	PermissionId *int64  `json:"permission_id"`
	Effect       *string `json:"effect"`
}
