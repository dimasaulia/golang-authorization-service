package dto

import "time"

type CreateUserPermissionOverrideRequest struct {
	UserId       int64      `json:"user_id"`
	PermissionId int64      `json:"permission_id"`
	Effect       string     `json:"effect"`
	Reason       string     `json:"reason"`
	ExpiresAt    *time.Time `json:"expires_at"`
	CreatedBy    *int64     `json:"created_by"`
}

type UpdateUserPermissionOverrideRequest struct {
	UserId       *int64     `json:"user_id"`
	PermissionId *int64     `json:"permission_id"`
	Effect       *string    `json:"effect"`
	Reason       *string    `json:"reason"`
	ExpiresAt    *time.Time `json:"expires_at"`
	CreatedBy    *int64     `json:"created_by"`
}
