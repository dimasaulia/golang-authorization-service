package dto

type CreateUserRequest struct {
	OrganizationId     int64   `json:"organization_id"`
	Username           string  `json:"username"`
	Email              string  `json:"email"`
	DisplayName        string  `json:"display_name"`
	Type               string  `json:"type"`
	Status             string  `json:"status"`
	Password           string  `json:"password"`
	MustChangePassword *bool   `json:"must_change_password"`
	CreateInKeycloak   bool    `json:"create_in_keycloak"`
	CreateInFreeIPA    bool    `json:"create_in_freeipa"`
	SendInvitation     bool    `json:"send_invitation"`
	RoleIds            []int64 `json:"role_ids"`
	TeamIds            []int64 `json:"team_ids"`
}

type UpdateUserRequest struct {
	OrganizationId *int64  `json:"organization_id"`
	Username       *string `json:"username"`
	Email          *string `json:"email"`
	DisplayName    *string `json:"display_name"`
	Type           *string `json:"type"`
	Status         *string `json:"status"`
}

type SignupUserRequest struct {
	OrganizationId int64  `json:"organization_id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	DisplayName    string `json:"display_name"`
	Password       string `json:"password"`
}

type GoogleSignupRequest struct {
	OrganizationId int64  `json:"organization_id"`
	ProviderUserId string `json:"provider_user_id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	DisplayName    string `json:"display_name"`
}

type VerifyEmailRequest struct {
	Code string `json:"code"`
}
