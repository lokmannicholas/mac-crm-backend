package entities

import (
	"context"

	"dmglab.com/mac-crm/pkg/models"
)

type Account struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Role        *Role  `json:"role"`
	Status      string `json:"status" enums:"Active,Disable" default:"Active"`
	LastLogin   *int64 `json:"last_login"`
}

func NewAccountEntity(account *models.Account) *Account {
	var lastLogin *int64
	if account == nil {
		return nil
	}
	if account.LastLogin.Valid {
		l := account.LastLogin.Time.Unix()
		lastLogin = &l
	}

	return &Account{
		ID:          account.ID.String(),
		DisplayName: account.DisplayName,
		Username:    account.Username,
		Role:        NewRoleEntity(&account.Role),
		Status:      account.Status,
		LastLogin:   lastLogin,
	}
}

func NewAccountListEntity(ctx context.Context, total int64, accounts []*models.Account) *List {
	accountList := make([]*Account, len(accounts))
	for i, account := range accounts {
		accountList[i] = NewAccountEntity(account)
	}
	return &List{
		Total: total,
		Data:  accountList,
	}

}

type Role struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Permissions      string `json:"permissions" example:"ACCOUNT:U;CUSTOMER:U;ROLE:U"`
	FieldPermissions string ` json:"field_permissions" example:"id_no;birth"`
	Levels           string ` json:"levels"`
}

func NewRoleEntity(role *models.Role) *Role {

	return &Role{
		ID:               role.ID.String(),
		Name:             role.Name,
		Permissions:      role.GetPermissions(),
		FieldPermissions: role.FieldPermissions,
		Levels:           role.Levels,
	}
}

func NewRoleListEntity(total int64, roles []*models.Role) *List {
	roleList := make([]*Role, len(roles))
	for i, role := range roles {
		roleList[i] = NewRoleEntity(role)
	}
	return &List{
		Total: total,
		Data:  roleList,
	}

}
