package dto

import "time"

type CreateUserRoleRequest struct {
	UserId         int64      `json:"user_id"`
	RoleId         int64      `json:"role_id"`
	AppId          *int64     `json:"app_id"`
	OrganizationId *int64     `json:"organization_id"`
	ExpiresAt      *time.Time `json:"expires_at"`
	AssignedBy     *int64     `json:"assigned_by"`
}

type UpdateUserRoleRequest struct {
	UserId         *int64     `json:"user_id"`
	RoleId         *int64     `json:"role_id"`
	AppId          *int64     `json:"app_id"`
	OrganizationId *int64     `json:"organization_id"`
	ExpiresAt      *time.Time `json:"expires_at"`
	AssignedBy     *int64     `json:"assigned_by"`
}
