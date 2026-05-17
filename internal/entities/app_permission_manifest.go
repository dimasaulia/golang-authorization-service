package entities

import (
	"encoding/json"
	"time"
)

type AppPermissionManifest struct {
	ID           int64            `db:"id" json:"id"`
	AppId        int64            `db:"app_id" json:"app_id"`
	Version      int64            `db:"version" json:"version"`
	Checksum     string           `db:"checksum" json:"checksum"`
	ManifestJson *json.RawMessage `db:"manifest_json" json:"manifest_json,omitempty"`
	CreatedAt    *time.Time       `db:"created_at" json:"created_at,omitempty"`
	ActivatedAt  *time.Time       `db:"activated_at" json:"activated_at,omitempty"`
}
