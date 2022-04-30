package models

import (
	"database/sql"

	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Feature struct {
	gorm.Model
	ID       uuid.UUID      `gorm:"type:varchar(36);primaryKey;"`
	Name     *MultiLangText `gorm:"embedded;embeddedPrefix:name_"`
	IsFilter sql.NullBool
	Parent   *Feature
	ParentID *uuid.UUID
	// BranchID uuid.UUID `gorm:"index:idx_branch_id"`
}

func (feature *Feature) BeforeCreate(tx *gorm.DB) (err error) {
	feature.ID = uuid.New()

	return
}

func (feature *Feature) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, feature)
	return
}

func (feature *Feature) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, feature)
	return
}
