package managers

import (
	"context"
	"errors"

	"dmglab.com/mac-crm/pkg/config"
	"dmglab.com/mac-crm/pkg/models"
	"dmglab.com/mac-crm/pkg/util"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SettingSetParam map[string]interface{}
type SettingGetParam struct {
	Key *string `json:"key"`
}
type ISettingManager interface {
	Set(ctx context.Context, param *SettingSetParam) ([]*models.App, error)
	GetByKey(ctx context.Context, param *SettingGetParam) (*models.App, error)
	GetAll(db *gorm.DB) ([]*models.App, error)
}

type SettingManager struct {
	config *config.Config
}

func GetSettingManager() ISettingManager {
	return &SettingManager{
		config: config.GetConfig(),
	}
}

func (m *SettingManager) Set(ctx context.Context, param *SettingSetParam) ([]*models.App, error) {

	if param == nil {
		return nil, errors.New("invalid input")
	}
	settings := make([]*models.App, len(*param))
	err := util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		i := 0
		for k, v := range *param {
			data := v.(string)
			settings[i] = &models.App{
				Setting: k,
				Value:   []byte(data),
			}
			i++
		}
		err := tx.Model(&settings).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "verion"}, {Name: "setting"}},
			DoUpdates: clause.AssignmentColumns([]string{"setting", "value"}),
		}).Create(&settings).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(settings) > 0 {
		conf := config.GetConfig()
		for _, setting := range settings {
			conf.Setting[setting.Setting] = string(setting.Value)
		}
	}
	return settings, err
}

func (m *SettingManager) GetByKey(ctx context.Context, param *SettingGetParam) (*models.App, error) {
	setting := new(models.App)
	return setting, util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		return tx.Model(setting).Where("setting = ? ", param.Key).First(setting).Error
	})
}

func (m *SettingManager) GetAll(db *gorm.DB) ([]*models.App, error) {
	companyField := []string{
		_const.IMAGE_COMPANY_NAME,
		_const.IMAGE_COMPANY_SHORT_NAME,
		_const.IMAGE_COMPANY_ADDRESS,
		_const.IMAGE_COMPANY_CODE,
		_const.IMAGE_COMPANY_PHONE,
		_const.IMAGE_COMPANY_EMAIL,
		_const.IMAGE_COMPANY_LOGO,
	}
	settings := []*models.App{}
	return settings, db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&settings).Where("setting not in (?)", companyField).Find(&settings).Error
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	})

}
