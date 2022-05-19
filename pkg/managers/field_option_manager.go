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

type IFieldOptionManager interface {
	Create(ctx context.Context, fieldID string, param *FieldOptionParam) (*models.FieldOption, error)
	Update(ctx context.Context, optionID string, param *FieldOptionParam) (*models.FieldOption, error)
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

func (m *FieldOptionManager) Create(ctx context.Context, fieldID string, param *FieldOptionParam) (*models.FieldOption, error) {
	parsedfieldID, err := uuid.Parse(fieldID)
	if err != nil {
		return nil, errors.New("Invalid field ID")
	}
	option := &models.FieldOption{
		ID:      uuid.New(),
		FieldID: &parsedfieldID,
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

func (m *FieldOptionManager) Update(ctx context.Context, optionID string, param *FieldOptionParam) (*models.FieldOption, error) {

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
