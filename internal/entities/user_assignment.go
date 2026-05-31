package entities

type UserAssignedRole struct {
	ID      int64   `db:"id" json:"id"`
	Code    string  `db:"code" json:"code"`
	Name    string  `db:"name" json:"name"`
	Scope   string  `db:"scope" json:"scope"`
	AppId   *int64  `db:"app_id" json:"app_id"`
	AppCode *string `db:"app_code" json:"app_code"`
	AppName *string `db:"app_name" json:"app_name"`
}

type UserAssignedTeam struct {
	ID   int64  `db:"id" json:"id"`
	Code string `db:"code" json:"code"`
	Name string `db:"name" json:"name"`
}
