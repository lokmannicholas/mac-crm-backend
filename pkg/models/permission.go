package models

import (
	"time"

	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	Name      string `gorm:"type:varchar(36);primaryKey;"`
	Read      bool
	Create    bool
	Update    bool
	Delete    bool
	BranchID  uuid.UUID `gorm:"type:varchar(36);primaryKey;"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time `gorm:"<-:update"`
}

func (permission *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	permission.CreatedAt = time.Now()
	return
}
func (permission *Permission) BeforeSave(tx *gorm.DB) (err error) {
	permission.UpdatedAt = time.Now()
	return
}

func (permission *Permission) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, permission)
	return
}

func (permission *Permission) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, permission)
	return
}
