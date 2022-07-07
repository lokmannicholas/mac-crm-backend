package models

import (
	"time"

	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	CreatedBy           *uuid.UUID       `json:"created_by"`
	UpdatedBy           *uuid.UUID       `json:"updated_by"`
	ID                  uuid.UUID        `gorm:"type:varchar(36);primaryKey;" json:"id"`
	Code                int64            `gorm:"not null;autoIncrement;unique;" json:"code"`
	FirstName           string           `json:"first_name"`
	LastName            string           `json:"last_name"`
	IDNo                string           `gorm:"unique;" json:"id_no"`
	Birth               *time.Time       `gorm:"type:DATETIME" json:"birth"`
	LoanDate            *time.Time       `gorm:"type:DATETIME" json:"loan_date"`
	CourtCaseFilingDate *time.Time       `gorm:"type:DATETIME" json:"court_case_filing_date"`
	CourtOrderDate      *time.Time       `gorm:"type:DATETIME" json:"court_order_date"`
	CourtReleaseDate    *time.Time       `gorm:"type:DATETIME" json:"court_release_date"`
	Status              string           `json:"status"`
	Levels              string           `json:"levels"`
	Meta                []*CustomersMeta `gorm:"foreignKey:CustomerID;references:ID" json:"meta"`
}

func (cus *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	cus.SetActive()
	return
}

func (cus *Customer) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, cus)
	return
}

func (cus *Customer) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, cus)
	return
}

func (cus *Customer) SetDisable() *Customer {
	cus.Status = "Disable"
	return cus
}
func (cus *Customer) SetActive() *Customer {
	cus.Status = "Active"
	return cus
}
