package models

import (
	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConsumableProductEvent struct {
	gorm.Model
	ID                  uuid.UUID `gorm:"type:varchar(36);primaryKey;"`
	Event               string
	ConsumableProductID uuid.UUID `gorm:"type:varchar(36);"`
	Update              []byte    `gorm:"type:text;"`
	TriggerByID         uuid.UUID `gorm:"type:varchar(36);"`
	TriggerBy           Account   `gorm:"<-:false;foreignKey:TriggerByID"`
}

func (productEvent *ConsumableProductEvent) BeforeCreate(tx *gorm.DB) (err error) {
	productEvent.ID = uuid.New()
	return
}
func (productEvent *ConsumableProductEvent) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, productEvent)
	return
}

func (productEvent *ConsumableProductEvent) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, productEvent)
	return
}
