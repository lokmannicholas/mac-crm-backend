package models

import (
	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomField struct {
	ID           uuid.UUID `gorm:"type:varchar(36);primaryKey;"`
	CustomObject string
	UniqueKey    string         `gorm:"type:varchar(18);unique;"`
	FieldName    *MultiLangText `gorm:"embedded;embeddedPrefix:field_name_"`
	Sorting      int32
	FieldType    string
	Remarks      string
	Status       string
	Options      []*FieldOption `gorm:"foreignKey:FieldID;references:ID" json:"field_options"`
}

func (cus *CustomField) BeforeCreate(tx *gorm.DB) (err error) {
	cus.Status = "Active"
	return
}

func (cus *CustomField) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, cus)
	return
}

func (cus *CustomField) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, cus)
	return
}
