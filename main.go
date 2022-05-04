package main

import (
	"fmt"

	"dmglab.com/mac-crm/pkg/api"
	"dmglab.com/mac-crm/pkg/collections"
	"dmglab.com/mac-crm/pkg/config"
	"dmglab.com/mac-crm/pkg/managers"
	"dmglab.com/mac-crm/pkg/service"
	"dmglab.com/mac-crm/pkg/util"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"gorm.io/gorm"
)

func init() {
	util.LoadENV()
}

func main() {
	conf := config.GetConfig()
	//get db information bylogin info
	service.SysLog.Infoln("initializing Mac ....")
	service.SysLog.Infoln(_const.INIT_LOGO)
	collections.GetCollection().Migrate(func(tx *gorm.DB) error {
		collections.Migration(tx)
		collections.InitalSystemAcc(tx)
		collections.InitalAccount(tx)
		return nil
	})

	settings, err := managers.GetSettingManager().GetAll(collections.GetCollection().DB)
	if err != nil {
		panic(err)
	}
	if len(settings) > 0 {
		for _, setting := range settings {
			conf.Setting[setting.Setting] = string(setting.Value)
		}
	}

	if err := api.GetRouter().Run(fmt.Sprintf(":%d", conf.App.Port)); err != nil {
		service.SysLog.Fatalln(err)
	}
}
