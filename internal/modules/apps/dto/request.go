package dto

type CreateAppRequest struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type UpdateAppRequest struct {
	Code   *string `json:"code"`
	Name   *string `json:"name"`
	Status *string `json:"status"`
}
