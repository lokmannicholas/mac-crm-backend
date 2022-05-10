package swagger

import "dmglab.com/mac-crm/pkg/entities"

type Account struct {
	Account *entities.Account `json:"account"`
}

type Accounts struct {
	Accounts *AccountList `json:"accounts"`
}

type AccountList struct {
	List
	Data []*entities.Account `json:"data"`
}

type ChangePasswordResponse struct {
	Token *Token `json:"token"`
}
