package models

import (
	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StorageEvent struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:varchar(36);primaryKey;"`
	Event       string
	StorageID   uuid.UUID `gorm:"type:varchar(36);"`
	Update      []byte    `gorm:"type:text;"`
	TriggerByID uuid.UUID `gorm:"type:varchar(36);"`
	TriggerBy   Account   `gorm:"<-:false;foreignKey:TriggerByID"`
}

func (storageEvent *StorageEvent) BeforeCreate(tx *gorm.DB) (err error) {
	storageEvent.ID = uuid.New()
	return
}
func (storageEvent *StorageEvent) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, storageEvent)
	return
}

func (storageEvent *StorageEvent) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, storageEvent)
	return
}
