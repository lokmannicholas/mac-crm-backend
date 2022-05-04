package models

import (
	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FieldOption struct {
	Name    *MultiLangText `gorm:"embedded;embeddedPrefix:name_" json:"name"`
	FieldID *uuid.UUID     `gorm:"type:varchar(36);" json:"field_id"`
}

func (option *FieldOption) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, option)
	return
}

func (option *FieldOption) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, option)
	return
}
