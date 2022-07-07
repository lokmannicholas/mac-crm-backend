package managers

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"dmglab.com/mac-crm/pkg/config"
	"dmglab.com/mac-crm/pkg/models"
	"dmglab.com/mac-crm/pkg/util"
	"github.com/google/uuid"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CustomerCreateParam struct {
	FirstName           string            `json:"first_name"`
	LastName            string            `json:"last_name"`
	IDNo                string            `json:"id_no"`
	Birth               *models.Date      `json:"birth"`
	LoanDate            *models.Date      `json:"loan_date"`
	CourtCaseFilingDate *models.Date      `json:"court_case_filing_date"`
	CourtOrderDate      *models.Date      `json:"court_order_date"`
	CourtReleaseDate    *models.Date      `json:"court_release_date"`
	Levels              []string          `json:"levels" example:"|1|2|"`
	Meta                map[string]string `json:"meta"`
}
type CustomerUpdateParam struct {
	FirstName           *string           `json:"first_name"`
	LastName            *string           `json:"last_name"`
	IDNo                *string           `json:"id_no"`
	Birth               *models.Date      `json:"birth"`
	LoanDate            *models.Date      `json:"loan_date"`
	CourtCaseFilingDate *models.Date      `json:"court_case_filing_date"`
	CourtOrderDate      *models.Date      `json:"court_order_date"`
	CourtReleaseDate    *models.Date      `json:"court_release_date"`
	Status              *string           `json:"status"`
	Levels              *[]string         `json:"levels" example:"|1|2|"`
	Meta                map[string]string `json:"meta"`
}

type CustomerQueryParam struct {
	FirstName         *string `form:"first_name" json:"first_name"`
	LastName          *string `form:"last_name" json:"last_name"`
	IDNo              *string `form:"id_no" json:"id_no"`
	CourtOrderDate    *string `form:"court_order_date" json:"court_order_date"`
	CourtReleaseDate  *string `form:"court_release_date" json:"court_release_date"`
	Email             *string `form:"email" json:"email"`
	Phone1            *string `form:"phone1" json:"phone1"`
	Phone2            *string `form:"phone2" json:"phone2"`
	Phone3            *string `form:"phone3" json:"phone3"`
	LoanType          *string `form:"loan_type" json:"loan_type"`
	IsBankrupt        *string `form:"is_bankrupt" json:"is_bankrupt"`
	IsDRP             *string `form:"is_drp" json:"is_drp"`
	IsIVA             *string `form:"is_iva" json:"is_iva"`
	CourtCaseInvolved *string `form:"court_case_involved" json:"court_case_involved"`
	Department        *string `form:"department" json:"department"`
}

type ICustomerManager interface {
	Create(ctx context.Context, accountID uuid.UUID, param *CustomerCreateParam) (*models.Customer, error)
	Update(ctx context.Context, accountID uuid.UUID, customerID string, param *CustomerUpdateParam) (*models.Customer, error)
	GetCustomers(ctx context.Context, param *CustomerQueryParam, fieldPermissions string, levels string) ([]*models.Customer, error)
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

func (m *CustomerManager) Create(ctx context.Context, accountID uuid.UUID, param *CustomerCreateParam) (*models.Customer, error) {
	cus := &models.Customer{
		CreatedBy: &accountID,
		ID:        uuid.New(),
		FirstName: param.FirstName,
		LastName:  param.LastName,
		IDNo:      param.IDNo,
	}
	if cus.Birth = nil; param.Birth != nil {
		cus.Birth = &param.Birth.Time
	}
	if cus.LoanDate = nil; param.LoanDate != nil {
		cus.LoanDate = &param.LoanDate.Time
	}
	if cus.CourtCaseFilingDate = nil; param.CourtCaseFilingDate != nil {
		cus.CourtCaseFilingDate = &param.CourtCaseFilingDate.Time
	}
	if cus.CourtOrderDate = nil; param.CourtOrderDate != nil {
		cus.CourtOrderDate = &param.CourtOrderDate.Time
	}
	if cus.CourtReleaseDate = nil; param.CourtReleaseDate != nil {
		cus.CourtReleaseDate = &param.CourtReleaseDate.Time
	}
	if cus.Levels = nil; param.Levels != nil {
		lvs, err := json.Marshal(param.Levels)
		if err != nil {
			return nil, err
		}
		levels := string(lvs)
		cus.Levels = &levels
	}

	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		// save meta
		for k, v := range param.Meta {
			meta := &models.CustomersMeta{
				Meta: &models.Meta{
					Key: k,
					Val: v,
				},
				CustomerID: cus.ID,
			}
			cus.Meta = append(cus.Meta, meta)
		}
		return tx.Create(cus).Error
	})
	if err != nil {
		return nil, err
	}
	return cus, nil
}
func (m *CustomerManager) Update(ctx context.Context, accountID uuid.UUID, customerID string, param *CustomerUpdateParam) (*models.Customer, error) {
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
		if param.IDNo != nil {
			cus.IDNo = *param.IDNo
		}
		if param.Status != nil {
			cus.Status = *param.Status
		}

		if cus.Birth = nil; param.Birth != nil {
			cus.Birth = &param.Birth.Time
		}
		if cus.LoanDate = nil; param.LoanDate != nil {
			cus.LoanDate = &param.LoanDate.Time
		}
		if cus.CourtCaseFilingDate = nil; param.CourtCaseFilingDate != nil {
			cus.CourtCaseFilingDate = &param.CourtCaseFilingDate.Time
		}
		if cus.CourtOrderDate = nil; param.CourtOrderDate != nil {
			cus.CourtOrderDate = &param.CourtOrderDate.Time
		}
		if cus.CourtReleaseDate = nil; param.CourtReleaseDate != nil {
			cus.CourtReleaseDate = &param.CourtReleaseDate.Time
		}
		if cus.Levels = nil; param.Levels != nil {
			lvs, err := json.Marshal(param.Levels)
			if err != nil {
				return err
			}
			levels := string(lvs)
			cus.Levels = &levels
		}

		err = tx.Save(cus).Error
		if err != nil {
			return err
		}

		// save meta
		if len(param.Meta) > 0 {
			for k, v := range param.Meta {
				meta := &models.CustomersMeta{
					Meta: &models.Meta{
						Key: k,
						Val: v,
					},
					CustomerID: cus.ID,
				}
				cus.Meta = append(cus.Meta, meta)

			}
			return tx.Model(cus.Meta).Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "customer_id"}, {Name: "key"}},
				DoUpdates: clause.AssignmentColumns([]string{"data_type", "val"})}).
				Create(cus.Meta).Error
		}
		return nil
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

func (m *CustomerManager) GetCustomers(ctx context.Context, param *CustomerQueryParam, fieldPermissions string, levels string) ([]*models.Customer, error) {
	cuss := []*models.Customer{}
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		var err error

		customerIds, paramCount, err := GetCustomerIdsByMetas(ctx, param)
		if err != nil {
			return err
		}

		if paramCount > 0 {
			tx = tx.Where("id IN ?", customerIds)
		}
		if param.IDNo != nil {
			keyword := "%" + *param.IDNo + "%"
			tx = tx.Where("id_no LIKE ?", keyword)
		}
		if param.FirstName != nil {
			keyword := "%" + *param.FirstName + "%"
			tx = tx.Where("first_name LIKE ?", keyword)
		}
		if param.LastName != nil {
			keyword := "%" + *param.LastName + "%"
			tx = tx.Where("last_name LIKE ?", keyword)
		}
		if param.Department != nil {
			keyword := "%" + *param.Department + "%"
			tx = tx.Where("levels LIKE ?", keyword)
		}
		if param.CourtOrderDate != nil {
			fromTime, toTime, err := util.StrToTimeRange(*param.CourtOrderDate)
			if err != nil {
				return err
			}
			tx = tx.Where("court_order_date BETWEEN ? AND ?", fromTime, toTime)
		}
		if param.CourtReleaseDate != nil {
			fromTime, toTime, err := util.StrToTimeRange(*param.CourtReleaseDate)
			if err != nil {
				return err
			}
			tx = tx.Where("court_release_date BETWEEN ? AND ?", fromTime, toTime)
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
		if len(levels) > 0 {
			userLevels := strings.Split(levels, ";")
			query := []string{}
			for _, v := range userLevels {
				query = append(query, fmt.Sprintf("JSON_CONTAINS(levels,'\"%s\"')", v))
			}
			if len(query) > 0 {
				tx = tx.Where("levels is null OR " + strings.Join(query, " OR "))
			} else {
				tx = tx.Where("levels is null")
			}
		}
		if param.Department != nil {
			userLevels := strings.Split(*param.Department, ",")
			query := []string{}
			for _, v := range userLevels {
				query = append(query, fmt.Sprintf("JSON_CONTAINS(levels,'\"%s\"')", v))
			}
			if len(query) > 0 {
				tx = tx.Where(strings.Join(query, " OR "))
			}
		}

		err = tx.Distinct().Find(&cuss).Error
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	})
	return cuss, err
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
		tx = tx.Where("id = ?", customerID)
		// check level permission
		if len(levels) > 0 {
			userLevels := strings.Split(levels, ";")
			query := []string{}
			for _, v := range userLevels {
				query = append(query, fmt.Sprintf("JSON_CONTAINS(levels,'\"%s\"')", v))
			}
			if len(query) > 0 {
				tx = tx.Where("levels is null OR " + strings.Join(query, " OR "))
			} else {
				tx = tx.Where("levels is null")
			}

		}
		return tx.First(cus).Error
	})

	return cus, err
}

func GetCustomerIdsByMetas(ctx context.Context, param *CustomerQueryParam) ([]uuid.UUID, int, error) {
	customersMetas := []*models.CustomersMeta{}
	paramCount := 0
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		if param.Email != nil {
			paramCount += 1
			keyword := "%" + *param.Email + "%"
			tx = tx.Or("`key` = 'email' AND val LIKE ?", keyword)
		}
		if param.Phone1 != nil {
			paramCount += 1
			keyword := "%" + *param.Phone1 + "%"
			tx = tx.Or("`key` = 'phone1' AND val LIKE ?", keyword)
		}
		if param.Phone2 != nil {
			paramCount += 1
			keyword := "%" + *param.Phone2 + "%"
			tx = tx.Or("`key` = 'phone2' AND val LIKE ?", keyword)
		}
		if param.Phone3 != nil {
			paramCount += 1
			keyword := "%" + *param.Phone3 + "%"
			tx = tx.Or("`key` = 'phone3' AND val LIKE ?", keyword)
		}
		if param.LoanType != nil {
			paramCount += 1
			keyword := "%" + *param.LoanType + "%"
			tx = tx.Or("`key` = 'loan_type' AND val LIKE ?", keyword)
		}
		if param.IsBankrupt != nil {
			paramCount += 1
			keyword := "%" + *param.IsBankrupt + "%"
			tx = tx.Or("`key` = 'is_bankrupt' AND val LIKE ?", keyword)
		}
		if param.IsDRP != nil {
			paramCount += 1
			keyword := "%" + *param.IsDRP + "%"
			tx = tx.Or("`key` = 'is_drp' AND val LIKE ?", keyword)
		}
		if param.IsIVA != nil {
			paramCount += 1
			keyword := "%" + *param.IsIVA + "%"
			tx = tx.Or("`key` = 'is_iva' AND val LIKE ?", keyword)
		}
		if param.CourtCaseInvolved != nil {
			paramCount += 1
			keyword := "%" + *param.CourtCaseInvolved + "%"
			tx = tx.Or("`key` = 'court_case_involved' AND val LIKE ?", keyword)
		}

		err := tx.Select("customer_id").Group("customer_id").Having("COUNT(*) >= ?", paramCount).Find(&customersMetas).Error
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	})

	customerIds := make([]uuid.UUID, 0)
	if err != nil {
		return customerIds, paramCount, err
	}

	for _, customersMeta := range customersMetas {
		customerIds = append(customerIds, customersMeta.CustomerID)
	}
	return customerIds, paramCount, err
}

func GetCustomerFields(ctx context.Context) []string {
	var customer interface{} = models.Customer{}
	val := reflect.ValueOf(customer)
	noOfFields := val.Type().NumField()
	fields := []string{}
	for i := 0; i < noOfFields; i++ {
		fieldName := val.Type().Field(i).Tag.Get("json")
		if fieldName != "meta" && fieldName != "" {
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
