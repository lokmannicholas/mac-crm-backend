package models

import (
	"time"

	_const "dmglab.com/mac-crm/pkg/util/const"

	"dmglab.com/mac-crm/pkg/service"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	ID            uuid.UUID `gorm:"type:varchar(36);primaryKey;"`
	PaymentNo     string    `gorm:"<-:create;unique;type:varchar(36)"`
	PaymentMethod string
	PaymentDate   *time.Time
	RefNo         string //for 3rd party payment reference
	TransactionID string
	RefFilename   string
	Ref           string //file path in server
	RefType       string
	Amount        float32
	Status        string
}

func (payment *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	payment.ID = uuid.New()
	return
}

func (payment *Payment) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, payment)
	return
}

func (payment *Payment) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, payment)
	return
}
