package util

import (
	"dmglab.com/mac-crm/pkg/service"
	"gorm.io/gorm"
)

func CheckAndCreate(migrator gorm.Migrator, v interface{}) error {
	if !migrator.HasTable(v) {
		err := migrator.CreateTable(v)
		if err != nil {
			service.SysLog.Errorln(err)
			return err
		}
	}
	return nil
}
