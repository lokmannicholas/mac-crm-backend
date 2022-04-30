package models

import (
	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:varchar(36);primaryKey;" json:"id"`
	OrderID     uuid.UUID `gorm:"type:varchar(36);primaryKey;" json:"order_id"`
	InvoiceNo   string    `gorm:"<-:create;unique;type:varchar(36)" json:"invoice_no"`
	InvoiceType string    `gorm:"type:varchar(36);" json:"invoice_type"`
	Amount      float32   `json:"amount"`
	PaymentID   uuid.UUID `gorm:"type:varchar(36);" json:"payment_id"`
	Payment     *Payment  `gorm:"<-:create;foreignKey:PaymentID" json:"-"`
}

func (invoice *Invoice) BeforeCreate(tx *gorm.DB) (err error) {
	invoice.ID = uuid.New()

	return
}
func (invoice *Invoice) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, invoice)
	return
}

func (invoice *Invoice) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, invoice)
	return
}
