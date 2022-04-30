package models

import (
	"time"

	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const NEW, CONTRACTED, INTERESTED, DESIRED, CLOSING = "NEW", "CONTRACTED", "INTERESTED", "DESIRED", "CLOSING"

type Lead struct {
	gorm.Model
	NewLeadDate      time.Time
	ContactedDate    *time.Time
	InterestedDate   *time.Time
	DesiredDate      *time.Time
	ClosingDate      *time.Time
	Status           string
	AssignToID       *uuid.UUID `gorm:"type:varchar(36);"`
	AssignTo         *Account   `gorm:"<-:false;foreignKey:AssignToID"`
	AssigneeID       uuid.UUID  `gorm:"type:varchar(36);"`
	Assignee         Account    `gorm:"<-:false;foreignKey:AssigneeID"`
	Priority         int
	LeadTitle        string
	LeadName         string
	LeadContact      string
	LeadFrom         string
	ShortDescription string
	Remarks          string
	CustomerID       *uuid.UUID  `gorm:"type:varchar(36);"`
	Customer         *Customer   `gorm:"<-:false;foreignKey:CustomerID"`
	LeadCalls        []*LeadCall `gorm:"<-:create;foreignKey:LeadID;references:ID"`
}

func (ld *Lead) BeforeCreate(tx *gorm.DB) (err error) {
	ld.Status = NEW
	ld.NewLeadDate = time.Now()
	ctx := tx.Statement.Context
	acc := ctx.Value("Account").(*Account)
	ld.AssigneeID = acc.ID
	return
}

func (ld *Lead) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, ld)
	return
}

func (ld *Lead) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, ld)
	return
}
