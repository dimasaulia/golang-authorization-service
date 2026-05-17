package dto

type CreateUserIdentityRequest struct {
	UserId         int64   `json:"user_id"`
	Provider       string  `json:"provider"`
	ProviderUserId string  `json:"provider_user_id"`
	Username       *string `json:"username"`
	Email          *string `json:"email"`
	IsPrimary      bool    `json:"is_primary"`
}

type UpdateUserIdentityRequest struct {
	UserId         *int64  `json:"user_id"`
	Provider       *string `json:"provider"`
	ProviderUserId *string `json:"provider_user_id"`
	Username       *string `json:"username"`
	Email          *string `json:"email"`
	IsPrimary      *bool   `json:"is_primary"`
}
