package managers

import (
	"context"

	"dmglab.com/mac-crm/pkg/config"
	"dmglab.com/mac-crm/pkg/models"
	"dmglab.com/mac-crm/pkg/util"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type CustomFieldCreateParam struct {
	CustomObject string              `json:"custom_object,omitempty"`
	Section      string              `json:"section,omitempty"`
	UniqueKey    string              `json:"unique_key,omitempty"`
	FieldName    *MultiLangText      `json:"field_name,omitempty"`
	FieldType    string              `json:"field_type,omitempty" enums:"string,multiple"`
	Remarks      string              `json:"remarks,omitempty"`
	Sorting      int32               `json:"sorting,omitempty"`
	Options      []*FieldOptionParam `json:"field_options,omitempty"`
}
type CustomFieldUpdateParam struct {
	FieldName *MultiLangText       `json:"field_name,omitempty"`
	Section   *string              `json:"section,omitempty"`
	UniqueKey *string              `json:"unique_key,omitempty"`
	FieldType *string              `json:"field_type,omitempty" enums:"string,multiple"`
	Remarks   *string              `json:"remarks,omitempty"`
	Sorting   *int32               `json:"sorting,omitempty"`
	Status    *string              `json:"status,omitempty" enums:"Active,Disable"`
	Options   *[]FieldOptionUpdate `json:"field_options,omitempty"`
}

type CustomFieldQueryParam struct {
	CustomObject *string `form:"custom_object" json:"custom_object,omitempty"`
	Status       *string `form:"status" json:"status,omitempty"`
	Page         int     `form:"page" json:"page"`
	Limit        int     `form:"limit" json:"limit"`
}

type ICustomFieldManager interface {
	Create(ctx context.Context, param *CustomFieldCreateParam) (*models.CustomField, error)
	Update(ctx context.Context, CustomFieldID string, param *CustomFieldUpdateParam) (*models.CustomField, error)
	GetCustomFields(ctx context.Context, param *CustomFieldQueryParam) ([]*models.CustomField, *util.Pagination, error)
}

type CustomFieldManager struct {
	config *config.Config
}

func GetCustomFieldManager() ICustomFieldManager {
	return &CustomFieldManager{
		config: config.GetConfig(),
	}
}

func (m *CustomFieldManager) Create(ctx context.Context, param *CustomFieldCreateParam) (*models.CustomField, error) {

	cus := &models.CustomField{
		ID:           uuid.New(),
		Section:      param.Section,
		CustomObject: param.CustomObject,
		FieldType:    param.FieldType,
		Remarks:      param.Remarks,
		UniqueKey:    param.UniqueKey,
		Sorting:      param.Sorting,
	}
	if param.FieldName != nil {
		cus.FieldName = &models.MultiLangText{
			En: param.FieldName.En,
			Zh: param.FieldName.Zh,
			Ch: param.FieldName.Ch,
		}
	}

	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		for _, v := range param.Options {
			if v != nil {
				option := &models.FieldOption{
					FieldID: &cus.ID,
				}
				if v.Name != nil {
					option.Name = &models.MultiLangText{
						En: v.Name.En,
						Zh: v.Name.Zh,
						Ch: v.Name.Ch,
					}
				}
				cus.Options = append(cus.Options, option)
			}
		}
		return tx.Create(cus).Error
	})
	if err != nil {
		return nil, err
	}
	return cus, nil
}
func (m *CustomFieldManager) Update(ctx context.Context, id string, param *CustomFieldUpdateParam) (*models.CustomField, error) {

	cus := new(models.CustomField)
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {

		err := tx.First(cus, "id = ?", id).Error
		if err != nil {
			return err
		}

		if param.FieldName != nil {
			cus.FieldName = &models.MultiLangText{
				En: param.FieldName.En,
				Zh: param.FieldName.Zh,
				Ch: param.FieldName.Ch,
			}
		}
		if param.Section != nil {
			cus.Section = *param.Section
		}
		if param.Sorting != nil {
			cus.Sorting = *param.Sorting
		}
		if param.FieldType != nil {
			cus.FieldType = *param.FieldType
		}
		if param.Remarks != nil {
			cus.Remarks = *param.Remarks
		}
		if param.Status != nil {
			cus.Status = *param.Status
		}
		err = tx.Save(cus).Error
		if err != nil {
			return err
		}

		// update option if not error
		if param.Options != nil {
			fieldOptMgr := GetFieldOptionManager()
			for _, v := range *param.Options {
				switch v.Action {
				case "CREATE":
					_, err := fieldOptMgr.Create(ctx, id, &FieldOptionParam{Name: v.Name})
					if err != nil {
						return err
					}
				case "UPDATE":
					_, err := fieldOptMgr.Update(ctx, v.ID, &FieldOptionParam{Name: v.Name})
					if err != nil {
						return err
					}
				case "DELETE":
					err := fieldOptMgr.Delete(ctx, v.ID)
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return GetCustomField(ctx, id)
}

func (m *CustomFieldManager) GetCustomFields(ctx context.Context, param *CustomFieldQueryParam) ([]*models.CustomField, *util.Pagination, error) {
	pagin := &util.Pagination{
		Limit: param.Limit,
		Page:  param.Page,
	}
	cuss := []*models.CustomField{}
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		if param.CustomObject != nil {
			tx = tx.Where("custom_object = ?", param.CustomObject)
		}
		if param.Status != nil {
			tx = tx.Where("status = ?", param.Status)
		}
		err := tx.Distinct().Scopes(util.PaginationScope(cuss, pagin, tx)).Preload("Options").Find(&cuss).Error
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	})
	return cuss, pagin, err
}

func GetCustomField(ctx context.Context, id string) (*models.CustomField, error) {
	customField := new(models.CustomField)
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		return tx.Preload("Options").First(customField, "id = ?", id).Error
	})

	return customField, err
}
