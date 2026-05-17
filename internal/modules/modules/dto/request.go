package dto

type CreateModuleRequest struct {
	AppId  int64  `json:"app_id"`
	Code   string `json:"code"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type UpdateModuleRequest struct {
	AppId  *int64  `json:"app_id"`
	Code   *string `json:"code"`
	Name   *string `json:"name"`
	Status *string `json:"status"`
}
