package swagger

import "dmglab.com/mac-crm/pkg/entities"

type Role struct {
	Role *entities.Role `json:"role"`
}

type Roles struct {
	Roles *RoleList `json:"roles"`
}

type RoleList struct {
	List
	Data []*entities.Role `json:"data"`
}
