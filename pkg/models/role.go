package models

import (
	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:varchar(36);primaryKey;"  json:"id"`
	Name        string    `gorm:"unique" json:"name"`
	Permissions string    ` json:"permissions"`
}

func (rol *Role) GetPermissions() string {
	if len(rol.Permissions) == 0 {
		return ""
	}
	return rol.Permissions
	// return encrypt.ASEDecrypt(rol.Permissions, config.GetConfig().ASEKey)
}

func (role *Role) BeforeCreate(tx *gorm.DB) (err error) {
	role.ID = uuid.New()
	return
}

func (role *Role) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, role)
	return
}

func (role *Role) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, role)
	return
}
