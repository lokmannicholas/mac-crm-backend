package models

type Auth struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	TenantID string `json:"tenant_id"`
}
