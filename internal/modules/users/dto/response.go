package dto

import "github.com/open-suite/authorization/internal/entities"

type ListResponse[T any] struct {
	Items []T `json:"items"`
}

type UserResponse struct {
	ID                 int64                       `json:"id"`
	OrganizationId     int64                       `json:"organization_id"`
	Username           string                      `json:"username"`
	Email              string                      `json:"email"`
	DisplayName        string                      `json:"display_name"`
	Type               string                      `json:"type"`
	Status             string                      `json:"status"`
	MustChangePassword bool                        `json:"must_change_password,omitempty"`
	ProvisionedTo      []string                    `json:"provisioned_to,omitempty"`
	RoleIds            []int64                     `json:"role_ids,omitempty"`
	TeamIds            []int64                     `json:"team_ids,omitempty"`
	Roles              []entities.UserAssignedRole `json:"roles,omitempty"`
	Teams              []entities.UserAssignedTeam `json:"teams,omitempty"`
}
