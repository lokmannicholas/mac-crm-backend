package managers

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"dmglab.com/mac-crm/pkg/config"
	"dmglab.com/mac-crm/pkg/models"
	"dmglab.com/mac-crm/pkg/util"
	"github.com/google/uuid"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CustomerCreateParam struct {
	FirstName            string                 `json:"first_name"`
	LastName             string                 `json:"last_name"`
	OtherName            string                 `json:"other_name"`
	Remarks              string                 `json:"remarks"`
	Code                 string                 `json:"code"`
	Adderess             string                 `json:"address"`
	Title                string                 `json:"title"`
	Phone1               string                 `json:"phone1"`
	Phone2               string                 `json:"phone2"`
	Phone3               string                 `json:"phone3"`
	IDNo                 string                 `gorm:"unique;" json:"id_no"`
	Birth                *time.Time             `json:"birth"`
	FbName               string                 `json:"fb_name"`
	IgName               string                 `json:"ig_name"`
	LinkedinName         string                 `json:"linkedin_name"`
	WeiboName            string                 `json:"weibo_name"`
	ResidentialAddress1  string                 `json:"residential_address_1"`
	ResidentialAddress2  string                 `json:"residential_address_2"`
	ResidentialAddress3  string                 `json:"residential_address_3"`
	OfficeAddress        string                 `json:"office_address"`
	AvgMonthIncome       string                 `json:"avg_month_income"`
	BankAcName           string                 `json:"bank_ac_name"`
	BankAcNumber         string                 `json:"bank_ac_number"`
	TaxIdNumber          string                 `json:"tax_id_number"`
	LoanAmount           string                 `json:"loan_amount"`
	LoanTenor            string                 `json:"loan_tenor"`
	LoanDate             string                 `json:"loan_date"`
	LoanInstalmentAmount string                 `json:"loan_instalment_amount"`
	DebtOutstanding      string                 `json:"debt_outstanding"`
	DebtInstalmentAmount string                 `json:"debt_instalment_amount"`
	DebtType             string                 `json:"debt_type"`
	DebtStatus           string                 `json:"debt_status"`
	Collateral           string                 `json:"collateral"`
	CollateralValue      string                 `json:"collateral_value"`
	Dti                  string                 `json:"dti"`
	Score                string                 `json:"score"`
	StatusDate           *time.Time             `json:"status_date"`
	OrderDate            *time.Time             `json:"order_date"`
	Levels               string                 `json:"levels" example:"|1|2|"`
	Meta                 map[string]interface{} `json:"meta"`
}
type CustomerUpdateParam struct {
	FirstName            *string                `json:"first_name"`
	LastName             *string                `json:"last_name"`
	OtherName            *string                `json:"other_name"`
	Remarks              *string                `json:"remarks"`
	Code                 *string                `json:"code"`
	Adderess             *string                `json:"address"`
	Title                *string                `json:"title"`
	Phone1               *string                `json:"phone1"`
	Phone2               *string                `json:"phone2"`
	Phone3               *string                `json:"phone3"`
	IDNo                 *string                `gorm:"unique;" json:"id_no"`
	Birth                *time.Time             `json:"birth"`
	FbName               *string                `json:"fb_name"`
	IgName               *string                `json:"ig_name"`
	LinkedinName         *string                `json:"linkedin_name"`
	WeiboName            *string                `json:"weibo_name"`
	ResidentialAddress1  *string                `json:"residential_address_1"`
	ResidentialAddress2  *string                `json:"residential_address_2"`
	ResidentialAddress3  *string                `json:"residential_address_3"`
	OfficeAddress        *string                `json:"office_address"`
	AvgMonthIncome       *string                `json:"avg_month_income"`
	BankAcName           *string                `json:"bank_ac_name"`
	BankAcNumber         *string                `json:"bank_ac_number"`
	TaxIdNumber          *string                `json:"tax_id_number"`
	LoanAmount           *string                `json:"loan_amount"`
	LoanTenor            *string                `json:"loan_tenor"`
	LoanDate             *string                `json:"loan_date"`
	LoanInstalmentAmount *string                `json:"loan_instalment_amount"`
	DebtOutstanding      *string                `json:"debt_outstanding"`
	DebtInstalmentAmount *string                `json:"debt_instalment_amount"`
	DebtType             *string                `json:"debt_type"`
	DebtStatus           *string                `json:"debt_status"`
	Collateral           *string                `json:"collateral"`
	CollateralValue      *string                `json:"collateral_value"`
	Dti                  *string                `json:"dti"`
	Score                *string                `json:"score"`
	StatusDate           *time.Time             `json:"status_date"`
	OrderDate            *time.Time             `json:"order_date"`
	Status               *string                `json:"status"`
	Levels               *string                `json:"levels" example:"|1|2|"`
	Meta                 map[string]interface{} `json:"meta"`
}

type CustomerQueryParam struct {
	SearchMode string  `form:"search_mode" json:"search_mode"`
	Code       *string `form:"code" json:"code"`
	Phone      *string `form:"phone" json:"phone"`
	IDNo       *string `form:"id_no" json:"id_no"`
	Page       int     `form:"page" json:"page"`
	Limit      int     `form:"limit" json:"limit"`
}

func (q *CustomerQueryParam) UnmarshalJSON(data []byte) error {

	var v map[string][]string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if len(v["search_mode"]) > 0 {
		q.SearchMode = v["search_mode"][0]
	}
	if len(v["code"]) > 0 {
		q.Code = &v["code"][0]
	}
	if len(v["phone"]) > 0 {
		q.Phone = &v["phone"][0]
	}
	if len(v["id_no"]) > 0 {
		q.IDNo = &v["id_no"][0]
	}
	if len(v["page"]) > 0 {
		i, err := strconv.ParseInt(v["page"][0], 10, 64)
		if err != nil {
			return err
		}
		q.Page = int(i)
	} else {
		q.Page = 0
	}
	return nil
}

type ICustomerManager interface {
	Create(ctx context.Context, param *CustomerCreateParam) (*models.Customer, error)
	Update(ctx context.Context, customerID string, param *CustomerUpdateParam) (*models.Customer, error)
	GetCustomers(ctx context.Context, param *CustomerQueryParam, fieldPermissions string, levels string) ([]*models.Customer, *util.Pagination, error)
	GetCustomer(ctx context.Context, customerID string, fieldPermissions string, levels string) (*models.Customer, error)
	Activate(ctx context.Context, customerID string) error
	Disable(ctx context.Context, customerID string) error
}

type CustomerManager struct {
	config *config.Config
}

func GetCustomerManager() ICustomerManager {
	return &CustomerManager{
		config: config.GetConfig(),
	}
}

func (m *CustomerManager) Create(ctx context.Context, param *CustomerCreateParam) (*models.Customer, error) {

	cus := &models.Customer{
		ID:                   uuid.New(),
		FirstName:            param.FirstName,
		LastName:             param.LastName,
		OtherName:            param.OtherName,
		Remarks:              param.Remarks,
		Code:                 param.Code,
		Adderess:             param.Adderess,
		Title:                param.Title,
		Phone1:               param.Phone1,
		Phone2:               param.Phone2,
		Phone3:               param.Phone3,
		IDNo:                 param.IDNo,
		Birth:                param.Birth,
		FbName:               param.FbName,
		IgName:               param.IgName,
		LinkedinName:         param.LinkedinName,
		WeiboName:            param.WeiboName,
		ResidentialAddress1:  param.ResidentialAddress1,
		ResidentialAddress2:  param.ResidentialAddress2,
		ResidentialAddress3:  param.ResidentialAddress3,
		OfficeAddress:        param.OfficeAddress,
		AvgMonthIncome:       param.AvgMonthIncome,
		BankAcName:           param.BankAcName,
		BankAcNumber:         param.BankAcNumber,
		TaxIdNumber:          param.TaxIdNumber,
		LoanAmount:           param.LoanAmount,
		LoanTenor:            param.LoanTenor,
		LoanDate:             param.LoanDate,
		LoanInstalmentAmount: param.LoanInstalmentAmount,
		DebtOutstanding:      param.DebtOutstanding,
		DebtInstalmentAmount: param.DebtInstalmentAmount,
		DebtType:             param.DebtType,
		DebtStatus:           param.DebtStatus,
		Collateral:           param.Collateral,
		CollateralValue:      param.CollateralValue,
		Dti:                  param.Dti,
		Score:                param.Score,
		StatusDate:           param.StatusDate,
		OrderDate:            param.OrderDate,
		Levels:               param.Levels,
	}
	if len(param.Code) == 0 {
		return nil, errors.New("customer code error")
	}
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		// save meta
		for k, v := range param.Meta {
			if v != nil {
				var buf bytes.Buffer
				enc := gob.NewEncoder(&buf)
				err := enc.Encode(v)
				if err != nil {
					return err
				}
				data := buf.Bytes()
				meta := &models.CustomersMeta{
					Meta: &models.Meta{
						Key: k,
						Val: data,
					},
					CustomerID: cus.ID,
				}
				cus.Meta = append(cus.Meta, meta)
			}
		}
		return tx.Create(cus).Error
	})
	if err != nil {
		return nil, err
	}
	return cus, nil
}
func (m *CustomerManager) Update(ctx context.Context, customerID string, param *CustomerUpdateParam) (*models.Customer, error) {

	if param.Code == nil {
		return nil, errors.New("customer code error")
	}
	cus := new(models.Customer)
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		err := tx.First(cus, "id = ?", customerID).Error
		if err != nil {
			return err
		}
		if param.FirstName != nil {
			cus.FirstName = *param.FirstName
		}
		if param.LastName != nil {
			cus.LastName = *param.LastName
		}
		if param.OtherName != nil {
			cus.OtherName = *param.OtherName
		}
		if param.IDNo != nil {
			cus.SetIDNo(*param.IDNo)
		}
		if param.Remarks != nil {
			cus.Remarks = *param.Remarks
		}
		if param.Code != nil {
			cus.Code = *param.Code
		}
		if param.Title != nil {
			cus.Title = *param.Title
		}
		if param.Phone1 != nil {
			cus.Phone1 = *param.Phone1
		}
		if param.Phone2 != nil {
			cus.Phone2 = *param.Phone2
		}
		if param.Phone3 != nil {
			cus.Phone3 = *param.Phone3
		}
		if param.IDNo != nil {
			cus.IDNo = *param.IDNo
		}
		if param.FbName != nil {
			cus.FbName = *param.FbName
		}
		if param.IgName != nil {
			cus.IgName = *param.IgName
		}
		if param.LinkedinName != nil {
			cus.LinkedinName = *param.LinkedinName
		}
		if param.WeiboName != nil {
			cus.WeiboName = *param.WeiboName
		}
		if param.ResidentialAddress1 != nil {
			cus.ResidentialAddress1 = *param.ResidentialAddress1
		}
		if param.ResidentialAddress2 != nil {
			cus.ResidentialAddress2 = *param.ResidentialAddress2
		}
		if param.ResidentialAddress3 != nil {
			cus.ResidentialAddress3 = *param.ResidentialAddress3
		}
		if param.OfficeAddress != nil {
			cus.OfficeAddress = *param.OfficeAddress
		}
		if param.AvgMonthIncome != nil {
			cus.AvgMonthIncome = *param.AvgMonthIncome
		}
		if param.BankAcName != nil {
			cus.BankAcName = *param.BankAcName
		}
		if param.BankAcNumber != nil {
			cus.BankAcNumber = *param.BankAcNumber
		}
		if param.TaxIdNumber != nil {
			cus.TaxIdNumber = *param.TaxIdNumber
		}
		if param.LoanAmount != nil {
			cus.LoanAmount = *param.LoanAmount
		}
		if param.LoanTenor != nil {
			cus.LoanTenor = *param.LoanTenor
		}
		if param.LoanDate != nil {
			cus.LoanDate = *param.LoanDate
		}
		if param.LoanInstalmentAmount != nil {
			cus.LoanInstalmentAmount = *param.LoanInstalmentAmount
		}
		if param.DebtOutstanding != nil {
			cus.DebtOutstanding = *param.DebtOutstanding
		}
		if param.DebtInstalmentAmount != nil {
			cus.DebtInstalmentAmount = *param.DebtInstalmentAmount
		}
		if param.DebtType != nil {
			cus.DebtType = *param.DebtType
		}
		if param.DebtStatus != nil {
			cus.DebtStatus = *param.DebtStatus
		}
		if param.Collateral != nil {
			cus.Collateral = *param.Collateral
		}
		if param.CollateralValue != nil {
			cus.CollateralValue = *param.CollateralValue
		}
		if param.Dti != nil {
			cus.Dti = *param.Dti
		}
		if param.Score != nil {
			cus.Score = *param.Score
		}
		if param.Adderess != nil {
			cus.Adderess = *param.Adderess
		}
		if param.Status != nil {
			cus.Status = *param.Status
		}
		if param.Levels != nil {
			cus.Levels = *param.Levels
		}
		cus.Birth = param.Birth
		cus.StatusDate = param.StatusDate
		cus.OrderDate = param.OrderDate

		err = tx.Save(cus).Error
		if err != nil {
			return err
		}

		// save meta
		for k, v := range param.Meta {
			if v != nil {
				var buf bytes.Buffer
				enc := gob.NewEncoder(&buf)
				err := enc.Encode(v)
				if err != nil {
					return err
				}
				data := buf.Bytes()
				meta := &models.CustomersMeta{
					Meta: &models.Meta{
						Key: k,
						Val: data,
					},
					CustomerID: cus.ID,
				}
				cus.Meta = append(cus.Meta, meta)
			}
		}
		return tx.Model(cus.Meta).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "customer_id"}, {Name: "key"}},
			DoUpdates: clause.AssignmentColumns([]string{"data_type", "val"})}).
			Create(cus.Meta).Error
	})

	return cus, err
}

func (m *CustomerManager) Disable(ctx context.Context, customerID string) error {

	return util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		cus := new(models.Customer)
		err := tx.First(cus, "id = ?", customerID).Error
		if err != nil {
			return err
		}
		cus.SetDisable()
		return tx.Save(cus).Error
	})
}
func (m *CustomerManager) Activate(ctx context.Context, customerID string) error {
	return util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		cus := new(models.Customer)
		err := tx.First(cus, "id = ?", customerID).Error
		if err != nil {
			return err
		}
		cus.SetActive()
		return tx.Save(cus).Error
	})
}

func (m *CustomerManager) GetCustomers(ctx context.Context, param *CustomerQueryParam, fieldPermissions string, levels string) ([]*models.Customer, *util.Pagination, error) {
	pagin := &util.Pagination{
		Limit: param.Limit,
		Page:  param.Page,
	}
	cuss := []*models.Customer{}
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		var err error
		whereExist := false
		if param.SearchMode == "eq" {
			if param.Phone != nil {
				tx = tx.Where("phone = ?", *param.Phone)
			}
			if param.IDNo != nil {
				tx = tx.Where("id_no = ?", *param.IDNo)
			}
			if param.Code != nil {
				tx = tx.Where("code = ?", *param.Code)
			}
			if param.Phone != nil {
				tx = tx.Where("phone = ?", *param.Phone)
			}
			if param.IDNo != nil {
				tx = tx.Where("id_no = ?", *param.IDNo)
			}
		} else {
			if param.Phone != nil {
				p := fmt.Sprintf("(.*%s)", *param.Phone)
				if !whereExist {
					tx = tx.Where("phone RLIKE ?", p)
					whereExist = true
				} else {
					tx = tx.Or("phone RLIKE ?", p)
				}
			}
			if param.IDNo != nil {
				p := fmt.Sprintf("(.*%s)", *param.IDNo)
				if !whereExist {
					tx = tx.Where("id_no RLIKE ?", p)
					whereExist = true
				} else {
					tx = tx.Or("id_no RLIKE ?", p)
				}
			}
		}
		if param.Code != nil {
			p := fmt.Sprintf("(.*%s)", *param.Code)
			if !whereExist {
				tx = tx.Where("code RLIKE ?", p)
				whereExist = true
			} else {
				tx = tx.Or("code RLIKE ?", p)
			}
		}
		if param.Phone != nil {
			p := fmt.Sprintf("(.*%s)", *param.Phone)
			if !whereExist {
				tx = tx.Where("phone RLIKE ?", p)
				whereExist = true
			} else {
				tx = tx.Or("phone RLIKE ?", p)
			}
		}
		if param.IDNo != nil {
			p := fmt.Sprintf("(.*%s)", *param.IDNo)
			if !whereExist {
				tx = tx.Where("id_no RLIKE ?", p)
				whereExist = true
			} else {
				tx = tx.Or("id_no RLIKE ?", p)
			}

		}

		// check field permission
		fields := strings.Split(fieldPermissions, ";")
		customerFields, metaFields := []string{"id"}, []string{}
		if fieldPermissions == "*" {
			tx = tx.Preload("Meta")
		} else {
			customerFieldsAllow, metaFieldsAllow := GetCustomerFields(ctx), GetMetaFields(ctx)
			for _, field := range fields {
				if util.Contain(customerFieldsAllow, field) {
					customerFields = append(customerFields, field)
				}
				if util.Contain(metaFieldsAllow, field) {
					metaFields = append(metaFields, field)
				}
			}
			tx = tx.Preload("Meta", "`key` IN ?", metaFields).Select(customerFields)
		}

		// check level permission
		userLevels := strings.Split(levels, ";")
		if levels != "*" {
			query := "levels = ''"
			for _, v := range userLevels {
				query = query + " OR levels LIKE '%|" + v + "|%'"
			}
			tx = tx.Where(query)
		}

		err = tx.Distinct().Scopes(util.PaginationScope(cuss, pagin, tx)).Find(&cuss).Error
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	})
	return cuss, pagin, err
}

func (m *CustomerManager) GetCustomer(ctx context.Context, customerID string, fieldPermissions string, levels string) (*models.Customer, error) {
	cus := new(models.Customer)
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		// check field permission
		fields := strings.Split(fieldPermissions, ";")
		customerFields, metaFields := []string{"id"}, []string{}
		if fieldPermissions == "*" {
			tx = tx.Preload("Meta")
		} else {
			customerFieldsAllow, metaFieldsAllow := GetCustomerFields(ctx), GetMetaFields(ctx)
			for _, field := range fields {
				if util.Contain(customerFieldsAllow, field) {
					customerFields = append(customerFields, field)
				}
				if util.Contain(metaFieldsAllow, field) {
					metaFields = append(metaFields, field)
				}
			}
			tx = tx.Preload("Meta", "`key` IN ?", metaFields).Select(customerFields)
		}

		// check level permission
		userLevels := strings.Split(levels, ";")
		if levels != "*" {
			query := "levels = ''"
			for _, v := range userLevels {
				query = query + " OR levels LIKE '%|" + v + "|%'"
			}
			tx = tx.Where(query)
		}

		return tx.Where("id = ?", customerID).First(cus).Error
	})

	return cus, err
}

func GetCustomerFields(ctx context.Context) []string {
	var customer interface{} = models.Customer{}
	val := reflect.ValueOf(customer)
	noOfFields := val.Type().NumField()
	fields := []string{}
	for i := 0; i < noOfFields; i++ {
		fieldName := val.Type().Field(i).Tag.Get("json")
		if fieldName != "meta" {
			fields = append(fields, fieldName)
		}
	}
	return fields
}

func GetMetaFields(ctx context.Context) []string {

	fields := []string{}
	customFields := []*models.CustomField{}
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		err := tx.Find(&customFields).Error
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	})
	if err == nil {
		for _, meta := range customFields {
			fields = append(fields, meta.UniqueKey)
		}
	}
	return fields
}
