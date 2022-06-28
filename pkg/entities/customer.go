package entities

import (
	"context"
	"strings"
	"time"

	"dmglab.com/mac-crm/pkg/managers"
	"dmglab.com/mac-crm/pkg/models"
	"github.com/google/uuid"
)

type Customer struct {
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	CreatedBy           *uuid.UUID `json:"created_by"`
	UpdatedBy           *uuid.UUID `json:"updated_by"`
	ID                  string     `json:"id"`
	FirstName           string     `json:"first_name"`
	LastName            string     `json:"last_name"`
	IDNo                string     `json:"id_no"`
	Birth               *time.Time `json:"birth"`
	LoanDate            *time.Time `json:"loan_date"`
	CourtCaseFilingDate *time.Time `json:"court_case_filing_date"`
	CourtOrderDate      *time.Time `json:"court_order_date"`
	CourtReleaseDate    *time.Time `json:"court_release_date"`
	Age                 *int       `json:"age"`
	Status              string     `json:"status" enums:"Active,Disable" default:"Active"`
	Levels              string     `json:"levels"`
	Meta                []Meta     `json:"meta"`
}

func NewCustomerEntity(customer *models.Customer, ctx context.Context) *Customer {
	if customer == nil {
		return &Customer{}
	}
	var age *int = nil
	if customer.Birth != nil {
		ageRes := calAge(*customer.Birth)
		age = &ageRes
	}
	metaArray := make([]Meta, len(customer.Meta))
	for i, meta := range customer.Meta {
		ent := *NewMetaEntity(meta.Meta)

		if ent.DataType == "multiple" {
			optionIds := strings.Split(meta.Val, ";")
			options, err := managers.FindByIds(ctx, optionIds)
			if err == nil {
				ent.Val = NewFieldOptinoListEntity(options)
			}
		}

		metaArray[i] = ent
	}
	return &Customer{
		CreatedAt:           customer.CreatedAt,
		UpdatedAt:           customer.UpdatedAt,
		CreatedBy:           customer.CreatedBy,
		UpdatedBy:           customer.UpdatedBy,
		ID:                  customer.ID.String(),
		FirstName:           customer.FirstName,
		LastName:            customer.LastName,
		IDNo:                customer.IDNo,
		Birth:               customer.Birth,
		CourtCaseFilingDate: customer.CourtCaseFilingDate,
		CourtOrderDate:      customer.CourtOrderDate,
		CourtReleaseDate:    customer.CourtReleaseDate,
		LoanDate:            customer.LoanDate,
		Age:                 age,
		Status:              customer.Status,
		Levels:              customer.Levels,
		Meta:                metaArray,
	}
}

func NewCustomerListEntity(customers []*models.Customer, ctx context.Context) *List {
	customerList := make([]*Customer, len(customers))
	columns := managers.GetCustomerFields(ctx)
	columns = append(columns, managers.GetMetaFields(ctx)...)
	for i, customer := range customers {
		customerList[i] = NewCustomerEntity(customer, ctx)
	}
	return &List{
		Columns: columns,
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
