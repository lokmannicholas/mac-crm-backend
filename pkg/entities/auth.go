package entities

import (
	"dmglab.com/mac-crm/pkg/models"
)

type Auth struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	TenantID string `json:"tenant_id"`
}

func NewAuthEntity(auth *models.Auth) *Auth {

	return &Auth{
		ID:       auth.ID,
		Username: auth.Username,
		TenantID: auth.TenantID,
	}
}
