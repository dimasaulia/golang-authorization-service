package dto

type CreateAppClientRequest struct {
	AppId            int64  `json:"app_id"`
	KeycloakClientId string `json:"keycloak_client_id"`
	Name             string `json:"name"`
	Status           string `json:"status"`
}

type UpdateAppClientRequest struct {
	AppId            *int64  `json:"app_id"`
	KeycloakClientId *string `json:"keycloak_client_id"`
	Name             *string `json:"name"`
	Status           *string `json:"status"`
}
