package entities

import (
	"dmglab.com/mac-crm/pkg/models"
)

type Auth struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	TenantID string `json:"tenant_id"`
}

type LoginResponse struct {
	AuthUser *AuthUser `json:"authUser"`
	Token    *Token    `json:"token"`
}

type AuthUser struct {
	CompanyID string `json:"companyID"`
	Role      *Role  `json:"role"`
	Username  string `json:"username"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func NewAuthEntity(auth *models.Auth) *Auth {

	return &Auth{
		ID:       auth.ID,
		Username: auth.Username,
		TenantID: auth.TenantID,
	}
}
