package dto

type CreateUserRequest struct {
	OrganizationId int64  `json:"organization_id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	DisplayName    string `json:"display_name"`
	Type           string `json:"type"`
	Status         string `json:"status"`
}

type UpdateUserRequest struct {
	OrganizationId *int64  `json:"organization_id"`
	Username       *string `json:"username"`
	Email          *string `json:"email"`
	DisplayName    *string `json:"display_name"`
	Type           *string `json:"type"`
	Status         *string `json:"status"`
}
