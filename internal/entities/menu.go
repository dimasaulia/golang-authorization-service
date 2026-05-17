package entities

type Menu struct {
	ID                   int64  `db:"id" json:"id"`
	AppId                int64  `db:"app_id" json:"app_id"`
	ModuleId             int64  `db:"module_id" json:"module_id"`
	ParentId             *int64 `db:"parent_id" json:"parent_id,omitempty"`
	Code                 string `db:"code" json:"code"`
	Name                 string `db:"name" json:"name"`
	RoutePath            string `db:"route_path" json:"route_path"`
	SortOrder            int64  `db:"sort_order" json:"sort_order"`
	RequiredPermissionId *int64 `db:"required_permission_id" json:"required_permission_id,omitempty"`
	Status               string `db:"status" json:"status"`
}
