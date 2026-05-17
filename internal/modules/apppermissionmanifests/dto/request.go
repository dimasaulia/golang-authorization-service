package dto

import (
	"encoding/json"
	"time"
)

type CreateAppPermissionManifestRequest struct {
	AppId        int64            `json:"app_id"`
	Version      int64            `json:"version"`
	Checksum     string           `json:"checksum"`
	ManifestJson *json.RawMessage `json:"manifest_json"`
	ActivatedAt  *time.Time       `json:"activated_at"`
}

type UpdateAppPermissionManifestRequest struct {
	AppId        *int64           `json:"app_id"`
	Version      *int64           `json:"version"`
	Checksum     *string          `json:"checksum"`
	ManifestJson *json.RawMessage `json:"manifest_json"`
	ActivatedAt  *time.Time       `json:"activated_at"`
}
