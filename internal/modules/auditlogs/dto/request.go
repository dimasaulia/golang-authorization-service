package dto

import "encoding/json"

type CreateAuditLogRequest struct {
	OrganizationId *int64           `json:"organization_id"`
	AppId          *int64           `json:"app_id"`
	ActorUserId    *int64           `json:"actor_user_id"`
	TargetUserId   *int64           `json:"target_user_id"`
	PermissionId   *int64           `json:"permission_id"`
	Action         string           `json:"action"`
	ResourceType   string           `json:"resource_type"`
	ResourceId     *string          `json:"resource_id"`
	Result         string           `json:"result"`
	MetadataJson   *json.RawMessage `json:"metadata_json"`
	IpAddress      *string          `json:"ip_address"`
	UserAgent      *string          `json:"user_agent"`
}

type UpdateAuditLogRequest struct {
	OrganizationId *int64           `json:"organization_id"`
	AppId          *int64           `json:"app_id"`
	ActorUserId    *int64           `json:"actor_user_id"`
	TargetUserId   *int64           `json:"target_user_id"`
	PermissionId   *int64           `json:"permission_id"`
	Action         *string          `json:"action"`
	ResourceType   *string          `json:"resource_type"`
	ResourceId     *string          `json:"resource_id"`
	Result         *string          `json:"result"`
	MetadataJson   *json.RawMessage `json:"metadata_json"`
	IpAddress      *string          `json:"ip_address"`
	UserAgent      *string          `json:"user_agent"`
}
