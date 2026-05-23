package dto

type CreateRolePermissionRequest struct {
	RoleId       int64  `json:"role_id"`
	PermissionId int64  `json:"permission_id"`
	Effect       string `json:"effect"`
}

type CreateBulkRolePermissionRequest struct {
	RoleId        int64   `json:"role_id"`
	PermissionIds []int64 `json:"permission_id"`
	Effect        string  `json:"effect"`
}

type UpdateRolePermissionRequest struct {
	RoleId       *int64  `json:"role_id"`
	PermissionId *int64  `json:"permission_id"`
	Effect       *string `json:"effect"`
}

type UpdateRolePermissionByRoleRequest struct {
	PermissionIds []int64 `json:"permission_id"`
	Effect        string  `json:"effect"`
}
