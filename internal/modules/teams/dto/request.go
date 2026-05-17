package dto

type CreateTeamRequest struct {
	OrganizationId int64  `json:"organization_id"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	Status         string `json:"status"`
}

type UpdateTeamRequest struct {
	OrganizationId *int64  `json:"organization_id"`
	Code           *string `json:"code"`
	Name           *string `json:"name"`
	Status         *string `json:"status"`
}
