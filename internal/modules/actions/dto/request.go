package dto

type CreateActionRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type UpdateActionRequest struct {
	Code *string `json:"code"`
	Name *string `json:"name"`
}
