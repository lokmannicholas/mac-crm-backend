package managers

import (
	"context"
	"errors"

	"dmglab.com/mac-crm/pkg/config"
	"dmglab.com/mac-crm/pkg/models"
	"dmglab.com/mac-crm/pkg/util"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type FieldOptionCreateParam struct {
	FieldOption
	FieldID *string `json:"field_id,omitempty"`
}
type FieldOptionUpdateParam struct {
	FieldOption
}

type IFieldOptionManager interface {
	Create(ctx context.Context, param *FieldOptionCreateParam) (*models.FieldOption, error)
	Update(ctx context.Context, optionID string, param *FieldOptionUpdateParam) (*models.FieldOption, error)
	Delete(ctx context.Context, optionID string) error
}

type FieldOptionManager struct {
	config *config.Config
}

func GetFieldOptionManager() IFieldOptionManager {
	return &FieldOptionManager{
		config: config.GetConfig(),
	}
}

func (m *FieldOptionManager) Create(ctx context.Context, param *FieldOptionCreateParam) (*models.FieldOption, error) {
	if param.FieldID == nil {
		return nil, errors.New("field ID can't be null")
	}
	fieldID, err := uuid.Parse(*param.FieldID)
	if err != nil {
		return nil, errors.New("Invalid field ID")
	}
	option := &models.FieldOption{
		ID:      uuid.New(),
		FieldID: &fieldID,
	}
	if param.Name != nil {
		option.Name = &models.MultiLangText{
			En: param.Name.En,
			Zh: param.Name.Zh,
			Ch: param.Name.Ch,
		}
	}

	err = util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		return tx.Create(option).Error
	})
	if err != nil {
		return nil, err
	}
	return option, nil
}

func (m *FieldOptionManager) Update(ctx context.Context, optionID string, param *FieldOptionUpdateParam) (*models.FieldOption, error) {

	option := new(models.FieldOption)
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {

		err := tx.First(option, "id = ?", optionID).Error
		if err != nil {
			return err
		}
		if param.Name != nil {
			option.Name = &models.MultiLangText{
				En: param.Name.En,
				Zh: param.Name.Zh,
				Ch: param.Name.Ch,
			}
		}
		err = tx.Save(option).Error
		if err != nil {
			return err
		}
		return nil
	})

	return option, err
}

func (m *FieldOptionManager) Delete(ctx context.Context, optionID string) error {
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		option := new(models.FieldOption)
		return tx.Where("id = ?", optionID).Delete(&option).Error
	})
	return err
}
