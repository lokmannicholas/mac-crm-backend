package models

import (
	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomersMeta struct {
	*Meta
	CustomerID uuid.UUID `gorm:"type:varchar(36);primaryKey;"`
}

func (p CustomersMeta) TableName() string {
	return "customers_meta"
}

func (m *CustomersMeta) BeforeCreate(tx *gorm.DB) (err error) {
	c := new(CustomField)
	err = tx.Where("status = ?", "Active").Where("unique_key = ?", m.Key).First(c).Error
	if err != nil {
		return err
	}
	m.DataType = c.FieldType
	return
}

func (m *CustomersMeta) BeforeSave(tx *gorm.DB) (err error) {
	c := new(CustomField)
	err = tx.Where("status = ?", "Active").Where("unique_key = ?", m.Key).First(c).Error
	if err != nil {
		return err
	}
	m.DataType = c.FieldType
	return
}

func (m *CustomersMeta) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, m)
	return
}

func (m *CustomersMeta) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, m)
	return
}
