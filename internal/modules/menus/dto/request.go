package dto

type CreateMenuRequest struct {
	AppId                int64  `json:"app_id"`
	ModuleId             int64  `json:"module_id"`
	ParentId             *int64 `json:"parent_id"`
	Code                 string `json:"code"`
	Name                 string `json:"name"`
	RoutePath            string `json:"route_path"`
	SortOrder            int64  `json:"sort_order"`
	RequiredPermissionId *int64 `json:"required_permission_id"`
	Status               string `json:"status"`
}

type UpdateMenuRequest struct {
	AppId                *int64  `json:"app_id"`
	ModuleId             *int64  `json:"module_id"`
	ParentId             *int64  `json:"parent_id"`
	Code                 *string `json:"code"`
	Name                 *string `json:"name"`
	RoutePath            *string `json:"route_path"`
	SortOrder            *int64  `json:"sort_order"`
	RequiredPermissionId *int64  `json:"required_permission_id"`
	Status               *string `json:"status"`
}
