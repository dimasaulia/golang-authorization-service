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
