package entities

import (
	"encoding/json"
	"time"
)

type AuditLog struct {
	ID             int64            `db:"id" json:"id"`
	OrganizationId *int64           `db:"organization_id" json:"organization_id,omitempty"`
	AppId          *int64           `db:"app_id" json:"app_id,omitempty"`
	ActorUserId    *int64           `db:"actor_user_id" json:"actor_user_id,omitempty"`
	TargetUserId   *int64           `db:"target_user_id" json:"target_user_id,omitempty"`
	PermissionId   *int64           `db:"permission_id" json:"permission_id,omitempty"`
	Action         string           `db:"action" json:"action"`
	ResourceType   string           `db:"resource_type" json:"resource_type"`
	ResourceId     *string          `db:"resource_id" json:"resource_id,omitempty"`
	Result         string           `db:"result" json:"result"`
	MetadataJson   *json.RawMessage `db:"metadata_json" json:"metadata_json,omitempty"`
	IpAddress      *string          `db:"ip_address" json:"ip_address,omitempty"`
	UserAgent      *string          `db:"user_agent" json:"user_agent,omitempty"`
	CreatedAt      *time.Time       `db:"created_at" json:"created_at,omitempty"`
}
