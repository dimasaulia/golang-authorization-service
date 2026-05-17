package entities

type Module struct {
	ID     int64  `db:"id" json:"id"`
	AppId  int64  `db:"app_id" json:"app_id"`
	Code   string `db:"code" json:"code"`
	Name   string `db:"name" json:"name"`
	Status string `db:"status" json:"status"`
}
