package models

import (
	"time"

	_const "dmglab.com/mac-crm/pkg/util/const"

	"dmglab.com/mac-crm/pkg/service"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RentRecord struct {
	gorm.Model
	ID            uuid.UUID `gorm:"type:varchar(36);primaryKey;"`
	RentStartDate time.Time
	RentEndDate   time.Time
	RentType      string //month year day
	RentPeriod    int
	Price         float32
	//belongs to
	Customer   Customer  `gorm:"foreignKey:CustomerID"`
	CustomerID uuid.UUID `gorm:"type:varchar(36);"`

	StorageID uuid.UUID `gorm:"type:varchar(36);"`
	OrderID   uuid.UUID `gorm:"type:varchar(36);"`
	OrderNo   string
	OrderDate time.Time
}

func (rentRecord *RentRecord) BeforeCreate(tx *gorm.DB) (err error) {
	rentRecord.ID = uuid.New()
	return
}

func (rentRecord *RentRecord) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, rentRecord)
	return
}

func (rentRecord *RentRecord) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, rentRecord)
	return
}
