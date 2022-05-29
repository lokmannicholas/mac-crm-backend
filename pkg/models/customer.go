package models

import (
	"time"

	"dmglab.com/mac-crm/pkg/config"
	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"dmglab.com/mac-crm/pkg/util/encrypt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	CreatedBy            *uuid.UUID       `json:"created_by"`
	UpdatedBy            *uuid.UUID       `json:"updated_by"`
	ID                   uuid.UUID        `gorm:"type:varchar(36);primaryKey;" json:"id"`
	Code                 string           `gorm:"<-:create;unique;" json:"code"`
	Adderess             string           `json:"address"`
	FirstName            string           `json:"first_name"`
	LastName             string           `json:"last_name"`
	OtherName            string           `json:"other_name"`
	Title                string           `json:"title"`
	Phone1               string           `json:"phone1"`
	Phone2               string           `json:"phone2"`
	Phone3               string           `json:"phone3"`
	IDNo                 string           `gorm:"unique;" json:"id_no"`
	Birth                *time.Time       `json:"birth"`
	FbName               string           `json:"fb_name"`
	IgName               string           `json:"ig_name"`
	LinkedinName         string           `json:"linkedin_name"`
	WeiboName            string           `json:"weibo_name"`
	ResidentialAddress1  string           `json:"residential_address_1"`
	ResidentialAddress2  string           `json:"residential_address_2"`
	ResidentialAddress3  string           `json:"residential_address_3"`
	OfficeAddress        string           `json:"office_address"`
	AvgMonthIncome       string           `json:"avg_month_income"`
	BankAcName           string           `json:"bank_ac_name"`
	BankAcNumber         string           `json:"bank_ac_number"`
	TaxIdNumber          string           `json:"tax_id_number"`
	LoanAmount           string           `json:"loan_amount"`
	LoanTenor            string           `json:"loan_tenor"`
	LoanDate             string           `json:"loan_date"`
	LoanInstalmentAmount string           `json:"loan_instalment_amount"`
	DebtOutstanding      string           `json:"debt_outstanding"`
	DebtInstalmentAmount string           `json:"debt_instalment_amount"`
	DebtType             string           `json:"debt_type"`
	DebtStatus           string           `json:"debt_status"`
	Collateral           string           `json:"collateral"`
	CollateralValue      string           `json:"collateral_value"`
	Dti                  string           `json:"dti"`
	Score                string           `json:"score"`
	StatusDate           *time.Time       `json:"status_date"`
	OrderDate            *time.Time       `json:"order_date"`
	Status               string           `json:"status"`
	Remarks              string           `json:"remarks"`
	Levels               string           `json:"levels"`
	Meta                 []*CustomersMeta `gorm:"foreignKey:CustomerID;references:ID" json:"meta"`
}

func (cus *Customer) BeforeCreate(tx *gorm.DB) (err error) {
	cus.SetActive()
	cus.SetIDNo(cus.IDNo)
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

func (cus *Customer) SetIDNo(idNo string) *Customer {
	cus.IDNo = encrypt.ASEEncrypt(idNo, config.GetConfig().ASEKey)
	return cus
}
func (cus *Customer) GetIDNo() string {
	if len(cus.IDNo) == 0 {
		return ""
	}
	decrypted := encrypt.ASEDecrypt(cus.IDNo, config.GetConfig().ASEKey)
	if len(decrypted) == 0 {
		return ""
	}
	return decrypted
}
