package entities

import (
	"bytes"
	"context"
	"encoding/gob"
	"strings"
	"time"

	"dmglab.com/mac-crm/pkg/managers"
	"dmglab.com/mac-crm/pkg/models"
	"github.com/google/uuid"
)

type Customer struct {
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	CreatedBy            *uuid.UUID `json:"created_by"`
	UpdatedBy            *uuid.UUID `json:"updated_by"`
	ID                   string     `json:"id"`
	Code                 string     `json:"code"`
	FirstName            string     `json:"first_name"`
	LastName             string     `json:"last_name"`
	OtherName            string     `json:"other_name"`
	Title                string     `json:"title"`
	Address              string     `json:"address"`
	Phone1               string     `json:"phone1"`
	Phone2               string     `json:"phone2"`
	Phone3               string     `json:"phone3"`
	IDNo                 string     `gorm:"unique;" json:"id_no"`
	Birth                *time.Time `json:"birth"`
	Age                  *int       `json:"age"`
	FbName               string     `json:"fb_name"`
	IgName               string     `json:"ig_name"`
	LinkedinName         string     `json:"linkedin_name"`
	WeiboName            string     `json:"weibo_name"`
	ResidentialAddress1  string     `json:"residential_address_1"`
	ResidentialAddress2  string     `json:"residential_address_2"`
	ResidentialAddress3  string     `json:"residential_address_3"`
	OfficeAddress        string     `json:"office_address"`
	AvgMonthIncome       string     `json:"avg_month_income"`
	BankAcName           string     `json:"bank_ac_name"`
	BankAcNumber         string     `json:"bank_ac_number"`
	TaxIdNumber          string     `json:"tax_id_number"`
	LoanAmount           string     `json:"loan_amount"`
	LoanTenor            string     `json:"loan_tenor"`
	LoanDate             string     `json:"loan_date"`
	LoanInstalmentAmount string     `json:"loan_instalment_amount"`
	DebtOutstanding      string     `json:"debt_outstanding"`
	DebtInstalmentAmount string     `json:"debt_instalment_amount"`
	DebtType             string     `json:"debt_type"`
	DebtStatus           string     `json:"debt_status"`
	Collateral           string     `json:"collateral"`
	CollateralValue      string     `json:"collateral_value"`
	Dti                  string     `json:"dti"`
	Score                string     `json:"score"`
	StatusDate           *time.Time `json:"status_date"`
	OrderDate            *time.Time `json:"order_date"`
	Status               string     `json:"status" enums:"Active,Disable" default:"Active"`
	Remarks              string     `json:"remarks"`
	Levels               string     `json:"levels"`
	Meta                 []Meta     `json:"meta"`
}

func NewCustomerEntity(customer *models.Customer, ctx context.Context) *Customer {
	if customer == nil {
		return &Customer{}
	}

	idNo := "******"
	if len(customer.GetIDNo()) > 4 {
		idNo = customer.GetIDNo()[0:4] + "******"
	}
	var age *int = nil
	if customer.Birth != nil {
		ageRes := calAge(*customer.Birth)
		age = &ageRes
	}
	metaArray := make([]Meta, len(customer.Meta))
	for i, meta := range customer.Meta {
		buf := bytes.NewBuffer(meta.Val)
		dec := gob.NewDecoder(buf)
		v := ""
		ent := *NewMetaEntity(meta.Meta)
		ent.Val = ""
		if err := dec.Decode(&v); err == nil {
			if ent.DataType == "multiple" {
				optionIds := strings.Split(v, ";")
				options, err := managers.FindByIds(ctx, optionIds)
				if err == nil {
					ent.Val = NewFieldOptinoListEntity(options)
				}
			} else {
				ent.Val = v
			}
		}
		metaArray[i] = ent
	}
	return &Customer{
		CreatedAt:            customer.CreatedAt,
		UpdatedAt:            customer.UpdatedAt,
		CreatedBy:            customer.CreatedBy,
		UpdatedBy:            customer.UpdatedBy,
		ID:                   customer.ID.String(),
		Code:                 customer.Code,
		FirstName:            customer.FirstName,
		LastName:             customer.LastName,
		OtherName:            customer.OtherName,
		Title:                customer.Title,
		Address:              customer.Adderess,
		Phone1:               customer.Phone1,
		Phone2:               customer.Phone2,
		Phone3:               customer.Phone3,
		IDNo:                 idNo,
		Birth:                customer.Birth,
		Age:                  age,
		FbName:               customer.FbName,
		IgName:               customer.IgName,
		LinkedinName:         customer.LinkedinName,
		WeiboName:            customer.WeiboName,
		ResidentialAddress1:  customer.ResidentialAddress1,
		ResidentialAddress2:  customer.ResidentialAddress2,
		ResidentialAddress3:  customer.ResidentialAddress3,
		OfficeAddress:        customer.OfficeAddress,
		AvgMonthIncome:       customer.AvgMonthIncome,
		BankAcName:           customer.BankAcName,
		BankAcNumber:         customer.BankAcNumber,
		TaxIdNumber:          customer.TaxIdNumber,
		LoanAmount:           customer.LoanAmount,
		LoanTenor:            customer.LoanTenor,
		LoanDate:             customer.LoanDate,
		LoanInstalmentAmount: customer.LoanInstalmentAmount,
		DebtOutstanding:      customer.DebtOutstanding,
		DebtInstalmentAmount: customer.DebtInstalmentAmount,
		DebtType:             customer.DebtType,
		DebtStatus:           customer.DebtStatus,
		Collateral:           customer.Collateral,
		CollateralValue:      customer.CollateralValue,
		Dti:                  customer.Dti,
		Score:                customer.Score,
		StatusDate:           customer.StatusDate,
		OrderDate:            customer.OrderDate,
		Status:               customer.Status,
		Remarks:              customer.Remarks,
		Levels:               customer.Levels,
		Meta:                 metaArray,
	}
}

func NewCustomerListEntity(total int64, customers []*models.Customer, ctx context.Context) *List {
	customerList := make([]*Customer, len(customers))
	columns := managers.GetCustomerFields(ctx)
	columns = append(columns, managers.GetMetaFields(ctx)...)
	for i, customer := range customers {
		customerList[i] = NewCustomerEntity(customer, ctx)
	}
	return &List{
		Columns: columns,
		Total:   total,
		Data:    customerList,
	}

}

func calAge(birthdate time.Time) int {
	today := time.Now()
	today = today.In(birthdate.Location())
	ty, tm, td := today.Date()
	today = time.Date(ty, tm, td, 0, 0, 0, 0, time.UTC)
	by, bm, bd := birthdate.Date()
	birthdate = time.Date(by, bm, bd, 0, 0, 0, 0, time.UTC)
	if today.Before(birthdate) {
		return 0
	}
	age := ty - by
	anniversary := birthdate.AddDate(age, 0, 0)
	if anniversary.After(today) {
		age--
	}
	return age
}
