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
	Code                int64            `gorm:"unique;" json:"code"`
	CreatedBy           *uuid.UUID       `json:"created_by"`
	UpdatedBy           *uuid.UUID       `json:"updated_by"`
	ID                  uuid.UUID        `gorm:"type:varchar(36);primaryKey;" json:"id"`
	FirstName           string           `json:"first_name"`
	LastName            string           `json:"last_name"`
	IDNo                string           `gorm:"unique;" json:"id_no"`
	Birth               *time.Time       `gorm:"type:DATETIME" json:"birth"`
	LoanDate            *time.Time       `gorm:"type:DATETIME" json:"loan_date"`
	CourtCaseFilingDate *time.Time       `gorm:"type:DATETIME" json:"court_case_filing_date"`
	CourtOrderDate      *time.Time       `gorm:"type:DATETIME" json:"court_order_date"`
	CourtReleaseDate    *time.Time       `gorm:"type:DATETIME" json:"court_release_date"`
	Status              string           `json:"status"`
	Levels              *string          `gorm:"type:json" json:"levels"`
	Meta                []*CustomersMeta `gorm:"foreignKey:CustomerID;references:ID" json:"meta"`
}

func (cus *Customer) BeforeCreate(tx *gorm.DB) error {
	var count int64
	err := tx.Model(&Customer{}).Count(&count).Error
	if err != nil {
		return err
	}
	cus.Code = count + 1
	cus.SetActive()
	cus.CreatedAt = time.Now()
	cus.UpdatedAt = time.Now()
	return nil
}
func (cus *Customer) BeforeSave(tx *gorm.DB) (err error) {
	cus.UpdatedAt = time.Now()
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
