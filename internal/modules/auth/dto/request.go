package dto

type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	SetCookie bool   `json:"set_cookie"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
	SetCookie    bool   `json:"set_cookie"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type KeycloakExchangeRequest struct {
	Code string `json:"code"`
}

type UpdateUserRequest struct {
	Username    *string `json:"username"`
	Email       *string `json:"email"`
	DisplayName *string `json:"display_name"`
	Password    *string `json:"password"`
}
