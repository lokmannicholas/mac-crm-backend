package models

import (
	"database/sql"

	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID       uuid.UUID      `gorm:"type:varchar(36);primaryKey;" json:"id"`
	Name     *MultiLangText `gorm:"embedded;embeddedPrefix:name_" json:"name"`
	IsFilter sql.NullBool   `json:"-"`
	Parent   *Category      `json:"-"`
	ParentID *uuid.UUID     `json:"-"`
	//BranchID uuid.UUID      `gorm:"index:idx_branch_id"`

}

func (category *Category) BeforeCreate(tx *gorm.DB) (err error) {
	category.ID = uuid.New()

	return
}

func (category *Category) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, category)
	return
}

func (category *Category) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, category)
	return
}
