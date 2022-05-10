package swagger

import "dmglab.com/mac-crm/pkg/entities"

type APIResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type List struct {
	Columns []string `json:"columns"`
	Total   int64    `json:"total"`
}

type APIForbiddenError struct {
	Message string `json:"message"`
}

type APIInternalServerError struct {
	Detail  string `json:"detail"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type LoginResponse struct {
	AuthUser *AuthUser `json:"authUser"`
	Token    *Token    `json:"token"`
}

type AuthUser struct {
	CompanyID string         `json:"companyID"`
	Role      *entities.Role `json:"role"`
	Username  string         `json:"username"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Settings struct {
	Settings *SettingList `json:"settings"`
}

type SettingList struct {
	List
	Data []*entities.Setting `json:"data"`
}

type Permissions struct {
	Permissions []string `json:"permissions"`
}
