package models

import (
	"dmglab.com/mac-crm/pkg/config"
	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"dmglab.com/mac-crm/pkg/util/encrypt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	ID        uuid.UUID        `gorm:"type:varchar(36);primaryKey;" json:"id"`
	Code      string           `gorm:"<-:create;unique;" json:"code"`
	Adderess  string           `json:"address"`
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	OtherName string           `json:"other_name"`
	Title     string           `json:"title"`
	Phone     string           `json:"phone"`
	IDNo      string           `gorm:"unique;" json:"id_no"`
	Status    string           `json:"status"`
	Remarks   string           `json:"remarks"`
	Meta      []*CustomersMeta `gorm:"foreignKey:CustomerID;references:ID" json:"meta"`
}

func (cus *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	cus.ID = uuid.New()
	cus.SetActive()
	cus.SetIDNo(cus.IDNo)
	return
}

func (cus *Customer) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, cus)
	return
}

func (cus *Customer) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, cus)
	return
}

func (cus *Customer) SetDisable() *Customer {
	cus.Status = "Disable"
	return cus
}
func (cus *Customer) SetActive() *Customer {
	cus.Status = "Active"
	return cus
}

func (cus *Customer) SetIDNo(idNo string) *Customer {
	cus.IDNo = encrypt.ASEEncrypt(idNo, config.GetConfig().ASEKey)
	return cus
}
func (cus *Customer) GetIDNo() string {
	if len(cus.IDNo) == 0 {
		return ""
	}
	decrypted := encrypt.ASEDecrypt(cus.IDNo, config.GetConfig().ASEKey)
	if len(decrypted) == 0 {
		return ""
	}
	return decrypted
}
