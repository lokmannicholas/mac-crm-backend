package models

import (
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Contract struct {
	ID         uuid.UUID      `gorm:"type:varchar(36);primaryKey;" json:"id"`
	Customer   Customer       `gorm:"<-:false;foreignKey:CustomerID" json:"-"`
	CustomerID uuid.UUID      `gorm:"type:varchar(36);" json:"customer_id"`
	Orders     []*RentalOrder `gorm:"foreignKey:ContractID;references:ID" json:"-"`
	Storage    Storage        `gorm:"foreignKey:ContractID;references:ID" json:"-"`
	Status     string         `json:"status"` //COMPLETE ACTIVE
	Agreements []*Agreement   `gorm:"foreignKey:ContractID;references:ID" json:"agreements"`
}

func (contract *Contract) BeforeCreate(tx *gorm.DB) (err error) {
	contract.ID = uuid.New()
	contract.Status = _const.CONTRACT_ACTIVE
	return
}

type Agreement struct {
	AgreementID uuid.UUID `gorm:"type:varchar(36);primaryKey;" json:"id"`
	ContractID  uuid.UUID `gorm:"type:varchar(36);" json:"contract_id"`
}
