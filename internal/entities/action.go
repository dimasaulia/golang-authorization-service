package entities

type Action struct {
	ID   int64  `db:"id" json:"id"`
	Code string `db:"code" json:"code"`
	Name string `db:"name" json:"name"`
}
