package swagger

import "dmglab.com/mac-crm/pkg/entities"

type Customer struct {
	Customer *entities.Customer `json:"customer"`
}

type Customers struct {
	Customers *CustomerList `json:"customers"`
}

type CustomerList struct {
	List
	Data []*entities.Customer `json:"data"`
}
