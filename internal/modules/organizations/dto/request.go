package dto

type CreateOrganizationRequest struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

type UpdateOrganizationRequest struct {
	Code   *string `json:"code"`
	Name   *string `json:"name"`
	Type   *string `json:"type"`
	Status *string `json:"status"`
}
