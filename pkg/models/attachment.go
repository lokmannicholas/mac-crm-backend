package models

import (
	"io"

	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Attachment struct {
	gorm.Model
	ID            uuid.UUID `gorm:"type:varchar(36);primaryKey;"`
	FileName      string
	Size          int
	Path          string
	Status        string
	MimeType      string
	StorageDriver string
	Reader        io.Reader `gorm:"-"`
}

func (att *Attachment) BeforeCreate(tx *gorm.DB) (err error) {
	att.Status = _const.FILE_VALID
	return
}

func (att *Attachment) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, att)
	return
}

func (att *Attachment) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, att)
	return
}

func (att *Attachment) AfterDelete(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, att)
	return
}
