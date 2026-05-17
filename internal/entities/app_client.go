package entities

import "time"

type AppClient struct {
	ID               int64      `db:"id" json:"id"`
	AppId            int64      `db:"app_id" json:"app_id"`
	KeycloakClientId string     `db:"keycloak_client_id" json:"keycloak_client_id"`
	Name             string     `db:"name" json:"name"`
	Status           string     `db:"status" json:"status"`
	CreatedAt        *time.Time `db:"created_at" json:"created_at,omitempty"`
}
