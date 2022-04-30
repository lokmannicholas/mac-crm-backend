package models

import (
	"time"

	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LeadCall struct {
	gorm.Model
	LeadID      uint
	CallingTime time.Time
	CallingVia  string
	CallRemark  string
	CalledByID  uuid.UUID
	CalledBy    Account
}

func (ld *LeadCall) BeforeCreate(tx *gorm.DB) (err error) {

	ld.CallingTime = time.Now()
	return
}

func (ld *LeadCall) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, ld)
	return
}

func (ld *LeadCall) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, ld)
	return
}
