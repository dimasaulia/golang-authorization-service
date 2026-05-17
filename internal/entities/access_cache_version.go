package entities

import "time"

type AccessCacheVersion struct {
	ID        int64      `db:"id" json:"id"`
	UserId    int64      `db:"user_id" json:"user_id"`
	AppId     int64      `db:"app_id" json:"app_id"`
	Version   int64      `db:"version" json:"version"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
