package entities

import "time"

type RolePermissionDetail struct {
	ID             int64      `db:"id" json:"id"`
	RoleId         int64      `db:"role_id" json:"role_id"`
	PermissionId   int64      `db:"permission_id" json:"permission_id"`
	Effect         string     `db:"effect" json:"effect"`
	CreatedAt      *time.Time `db:"created_at" json:"created_at,omitempty"`
	PermissionCode string     `db:"permission_code" json:"permission_code"`
	PermissionName string     `db:"permission_name" json:"permission_name"`
	ModuleId       *int64     `db:"module_id" json:"module_id,omitempty"`
	ModuleCode     *string    `db:"module_code" json:"module_code,omitempty"`
	ModuleName     *string    `db:"module_name" json:"module_name,omitempty"`
	AppId          int64      `db:"app_id" json:"app_id"`
	AppCode        string     `db:"app_code" json:"app_code"`
	AppName        string     `db:"app_name" json:"app_name"`
}

type RolePermissionSummary struct {
	RoleId          int64   `db:"role_id" json:"role_id"`
	RoleCode        string  `db:"role_code" json:"role_code"`
	RoleName        string  `db:"role_name" json:"role_name"`
	RoleDescription *string `db:"role_description" json:"role_description,omitempty"`
	RoleScope       string  `db:"role_scope" json:"role_scope"`
	AppId           *int64  `db:"app_id" json:"app_id,omitempty"`
	AppCode         *string `db:"app_code" json:"app_code,omitempty"`
	AppName         *string `db:"app_name" json:"app_name,omitempty"`
	PermissionCount int64   `db:"permission_count" json:"permission_count"`
}

type AvailablePermissionModule struct {
	ModuleId    *int64       `json:"modul_id,omitempty"`
	ModuleName  *string      `json:"modul_name,omitempty"`
	ModuleCode  *string      `json:"modul_code,omitempty"`
	Permissions []Permission `json:"permissions"`
}

type AvailablePermissionRow struct {
	ID          int64      `db:"id"`
	AppId       int64      `db:"app_id"`
	ModuleId    *int64     `db:"module_id"`
	ActionId    int64      `db:"action_id"`
	Code        string     `db:"code"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	RiskLevel   string     `db:"risk_level"`
	IsSystem    bool       `db:"is_system"`
	Status      string     `db:"status"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	ModuleName  *string    `db:"module_name"`
	ModuleCode  *string    `db:"module_code"`
}
