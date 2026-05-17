package dto

type CreateRoleRequest struct {
	OrganizationId *int64  `json:"organization_id"`
	AppId          *int64  `json:"app_id"`
	Code           string  `json:"code"`
	Name           string  `json:"name"`
	Description    *string `json:"description"`
	Scope          string  `json:"scope"`
	IsSystem       bool    `json:"is_system"`
	Status         string  `json:"status"`
}

type UpdateRoleRequest struct {
	OrganizationId *int64  `json:"organization_id"`
	AppId          *int64  `json:"app_id"`
	Code           *string `json:"code"`
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	Scope          *string `json:"scope"`
	IsSystem       *bool   `json:"is_system"`
	Status         *string `json:"status"`
}
