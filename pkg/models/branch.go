package models

import (
	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Branch struct {
	gorm.Model
	ID            uuid.UUID         `gorm:"type:varchar(36);primaryKey;" json:"id"`
	Code          string            `json:"code"`
	Name          *MultiLangText    `gorm:"embedded;embeddedPrefix:name_" json:"name"`
	ShortName     *MultiLangText    `gorm:"embedded;embeddedPrefix:short_name_" json:"short_name"`
	Address       *MultiLangText    `gorm:"embedded;embeddedPrefix:address_" json:"address"`
	Email         string            `json:"email"`
	Phone         string            `json:"phone"`
	TotalStorages int64             `gorm:"-" json:"-"`
	Attr          *BranchRentalAttr `gorm:"foreignKey:BranchID;references:ID" json:"-"`
}

func (brh *Branch) BeforeCreate(tx *gorm.DB) (err error) {
	brh.ID = uuid.New()
	return
}
func (brh *Branch) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, brh)
	return
}

func (brh *Branch) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, brh)
	return
}

type BranchRentalAttr struct {
	BranchID     uuid.UUID `gorm:"type:varchar(36);primaryKey;" json:"id"`
	AllowYearly  *bool
	AllowMonthly *bool
	AllowDaily   *bool
}
