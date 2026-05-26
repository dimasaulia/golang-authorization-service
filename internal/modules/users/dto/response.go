package dto

type ListResponse[T any] struct {
	Items []T `json:"items"`
}

type UserProvisioningRequest struct {
	Keycloak bool `json:"keycloak"`
	FreeIPA  bool `json:"freeipa"`
}

type UserResponse struct {
	ID                 int64                   `json:"id"`
	OrganizationId     int64                   `json:"organization_id"`
	Username           string                  `json:"username"`
	Email              string                  `json:"email"`
	DisplayName        string                  `json:"display_name"`
	Type               string                  `json:"type"`
	Status             string                  `json:"status"`
	MustChangePassword bool                    `json:"must_change_password,omitempty"`
	Provisioning       UserProvisioningRequest `json:"provisioning"`
	RoleIds            []int64                 `json:"role_ids,omitempty"`
	TeamIds            []int64                 `json:"team_ids,omitempty"`
}
