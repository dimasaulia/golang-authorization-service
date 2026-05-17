package entities

import "time"

type Permission struct {
	ID          int64      `db:"id" json:"id"`
	AppId       int64      `db:"app_id" json:"app_id"`
	ModuleId    *int64     `db:"module_id" json:"module_id,omitempty"`
	ActionId    int64      `db:"action_id" json:"action_id"`
	Code        string     `db:"code" json:"code"`
	Name        string     `db:"name" json:"name"`
	Description *string    `db:"description" json:"description,omitempty"`
	RiskLevel   string     `db:"risk_level" json:"risk_level"`
	IsSystem    bool       `db:"is_system" json:"is_system"`
	Status      string     `db:"status" json:"status"`
	CreatedAt   *time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}
