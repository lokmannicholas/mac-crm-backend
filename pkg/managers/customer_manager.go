package managers

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"dmglab.com/mac-crm/pkg/config"
	"dmglab.com/mac-crm/pkg/models"
	"dmglab.com/mac-crm/pkg/util"
	"github.com/google/uuid"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CustomerCreateParam struct {
	FirstName string                 `json:"first_name"`
	LastName  string                 `json:"last_name"`
	OtherName string                 `json:"other_name"`
	Phone     string                 `json:"phone"`
	IDNo      string                 `json:"id_no"`
	Remarks   string                 `json:"remarks"`
	Code      string                 `json:"code"`
	Adderess  string                 `json:"address"`
	Title     string                 `json:"title"`
	Meta      map[string]interface{} `json:"meta"`
}
type CustomerUpdateParam struct {
	FirstName *string                `json:"first_name"`
	LastName  *string                `json:"last_name"`
	OtherName *string                `json:"other_name"`
	Phone     *string                `json:"phone"`
	IDNo      *string                `json:"id_no"`
	Remarks   *string                `json:"remarks"`
	Code      *string                `json:"code"`
	Adderess  *string                `json:"address"`
	Title     *string                `json:"title"`
	Meta      map[string]interface{} `json:"meta"`
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
	GetCustomers(ctx context.Context, param *CustomerQueryParam) ([]*models.Customer, *util.Pagination, error)
	GetCustomer(ctx context.Context, customerID string) (*models.Customer, error)
	GetCustomerStorages(ctx context.Context, customerID string) ([]*models.Storage, error)
	GetCustomerSingleProduct(ctx context.Context, customerID string) ([]*models.SingleProduct, error)
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
		FirstName: param.FirstName,
		LastName:  param.LastName,
		OtherName: param.OtherName,
		Phone:     param.Phone,
		IDNo:      param.IDNo,
		Remarks:   param.Remarks,
		Code:      param.Code,
		Adderess:  param.Adderess,
		Title:     param.Title,
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
		if param.Phone != nil {
			cus.Phone = *param.Phone
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
		if param.Adderess != nil {
			cus.Adderess = *param.Adderess
		}

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
func (m *CustomerManager) GetCustomerStorages(ctx context.Context, customerID string) ([]*models.Storage, error) {
	var storages []*models.Storage
	return storages, util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		type customerStorage struct {
			OrderID uuid.UUID
			OrderNo string
			Item    []byte
		}
		items := []*customerStorage{}
		err := tx.Select("rental_order_items.order_id, rental_orders.order_no, rental_order_items.item").Model(models.RentalOrderItem{}).
			Joins("JOIN rental_orders on rental_orders.id = rental_order_items.order_id").
			Where("rental_orders.customer_id = ?", customerID).Scan(&items).Error
		if err != nil {
			return err
		}
		storages = make([]*models.Storage, len(items))
		for i, item := range items {

			storage := new(models.Storage)
			err = json.Unmarshal(item.Item, storage)
			if err != nil {
				return err
			}
			storage.OrderID = &item.OrderID
			storage.CurrentOrder = &models.RentalOrder{
				ID:      item.OrderID,
				OrderNo: item.OrderNo,
			}
			storages[i] = storage
		}
		return nil
	})
}

func (m *CustomerManager) GetCustomerSingleProduct(ctx context.Context, customerID string) ([]*models.SingleProduct, error) {
	singleProducts := []*models.SingleProduct{}
	return singleProducts, util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		err := tx.Model(models.SingleProduct{}).
			Where("customer_id = ?", customerID).Find(&singleProducts).Error
		if err != nil {
			return err
		}
		return nil
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

func (m *CustomerManager) GetCustomers(ctx context.Context, param *CustomerQueryParam) ([]*models.Customer, *util.Pagination, error) {
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

		err = tx.Preload("Meta").Distinct().Scopes(util.PaginationScope(cuss, pagin, tx)).Find(&cuss).Error
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	})
	return cuss, pagin, err
}

func (m *CustomerManager) GetCustomer(ctx context.Context, customerID string) (*models.Customer, error) {
	cus := new(models.Customer)
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		return tx.Preload("Meta").First(cus, "id = ?", customerID).Error
	})

	return cus, err
}
