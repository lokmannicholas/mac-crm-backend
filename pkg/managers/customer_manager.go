package managers

import (
	"context"
	"reflect"
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
	Name             string            `json:"name"`
	IDNo             string            `json:"id_no"`
	Birth            *time.Time        `json:"birth"`
	CourtFilingDate  *time.Time        `json:"court_filing_date"`
	CourtOrderDate   *time.Time        `json:"court_order_date"`
	CourtReleaseDate *time.Time        `json:"court_release_date"`
	Levels           string            `json:"levels" example:"|1|2|"`
	Meta             map[string]string `json:"meta"`
}
type CustomerUpdateParam struct {
	Name             *string           `json:"name"`
	IDNo             *string           `json:"id_no"`
	Birth            *time.Time        `json:"birth"`
	CourtFilingDate  *time.Time        `json:"court_filing_date"`
	CourtOrderDate   *time.Time        `json:"court_order_date"`
	CourtReleaseDate *time.Time        `json:"court_release_date"`
	Status           *string           `json:"status"`
	Levels           *string           `json:"levels" example:"|1|2|"`
	Meta             map[string]string `json:"meta"`
}

type CustomerQueryParam struct {
	Name              *string `form:"name" json:"name"`
	IDNo              *string `form:"id_no" json:"id_no"`
	CourtOrderDate    *string `form:"court_order_date" json:"court_order_date" example:"2022-05-14T00:00:00.000Z#2022-07-14T00:00:00.000Z"`
	CourtReleaseDate  *string `form:"court_release_date" json:"court_release_date" example:"2022-05-14T00:00:00.000Z#2022-07-14T00:00:00.000Z"`
	Email             *string `form:"email" json:"email"`
	Phone             *string `form:"phone" json:"phone"`
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
		CreatedBy:        &accountID,
		ID:               uuid.New(),
		Name:             param.Name,
		IDNo:             param.IDNo,
		Birth:            param.Birth,
		CourtFilingDate:  param.CourtFilingDate,
		CourtOrderDate:   param.CourtOrderDate,
		CourtReleaseDate: param.CourtReleaseDate,
		Levels:           param.Levels,
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
		if param.Name != nil {
			cus.Name = *param.Name
		}
		if param.IDNo != nil {
			cus.IDNo = *param.IDNo
		}
		if param.Status != nil {
			cus.Status = *param.Status
		}
		if param.Levels != nil {
			cus.Levels = *param.Levels
		}

		cus.Birth = param.Birth
		cus.CourtFilingDate = param.CourtFilingDate
		cus.CourtOrderDate = param.CourtOrderDate
		cus.CourtReleaseDate = param.CourtReleaseDate

		err = tx.Save(cus).Error
		if err != nil {
			return err
		}

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

func (m *CustomerManager) GetCustomers(ctx context.Context, param *CustomerQueryParam, fieldPermissions string, levels string) ([]*models.Customer, error) {
	cuss := []*models.Customer{}
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		var err error

		customerIds, err := GetCustomerIdsByMetas(ctx, param)
		if err != nil {
			return err
		}
		tx = tx.Where("id IN ?", customerIds)

		if param.IDNo != nil {
			keyword := "%" + *param.IDNo + "%"
			tx = tx.Where("id_no LIKE ?", keyword)
		}
		if param.Name != nil {
			keyword := "%" + *param.Name + "%"
			tx = tx.Where("name LIKE ?", keyword)
		}
		if param.CourtOrderDate != nil {
			dateSplit := strings.Split(*param.CourtOrderDate, "#")
			tx = tx.Where("court_order_date BETWEEN ? AND ?", dateSplit[0], dateSplit[1])
		}
		if param.CourtReleaseDate != nil {
			dateSplit := strings.Split(*param.CourtReleaseDate, "#")
			tx = tx.Where("court_release_date BETWEEN ? AND ?", dateSplit[0], dateSplit[1])
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

func GetCustomerIdsByMetas(ctx context.Context, param *CustomerQueryParam) ([]uuid.UUID, error) {
	customersMetas := []*models.CustomersMeta{}
	paramCount := 0
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		if param.Email != nil {
			paramCount += 1
			keyword := "%" + *param.Email + "%"
			tx = tx.Or("`key` = 'email' AND val LIKE ?", keyword)
		}
		if param.Phone != nil {
			paramCount += 1
			keyword := "%" + *param.Phone + "%"
			tx = tx.Or("(`key` = 'phone1' OR `key` = 'phone2' OR `key` = 'phone3') AND val LIKE ?", keyword)
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
		if param.Department != nil {
			paramCount += 1
			keyword := "%" + *param.Department + "%"
			tx = tx.Or("`key` = 'department' AND val LIKE ?", keyword)
		}

		err := tx.Select("customer_id").Group("customer_id").Having("COUNT(*) >= ?", paramCount).Find(&customersMetas).Error
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	})

	customerIds := make([]uuid.UUID, 0)
	if err != nil {
		return customerIds, err
	}

	for _, customersMeta := range customersMetas {
		customerIds = append(customerIds, customersMeta.CustomerID)
	}
	return customerIds, err
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
