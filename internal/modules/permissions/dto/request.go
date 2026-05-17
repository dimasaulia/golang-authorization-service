package dto

type CreatePermissionRequest struct {
	AppId       int64   `json:"app_id"`
	ModuleId    *int64  `json:"module_id"`
	ActionId    int64   `json:"action_id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	RiskLevel   string  `json:"risk_level"`
	IsSystem    bool    `json:"is_system"`
	Status      string  `json:"status"`
}

type UpdatePermissionRequest struct {
	AppId       *int64  `json:"app_id"`
	ModuleId    *int64  `json:"module_id"`
	ActionId    *int64  `json:"action_id"`
	Code        *string `json:"code"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	RiskLevel   *string `json:"risk_level"`
	IsSystem    *bool   `json:"is_system"`
	Status      *string `json:"status"`
}
