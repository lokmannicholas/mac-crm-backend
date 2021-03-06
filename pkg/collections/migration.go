package collections

import (
	"flag"

	"dmglab.com/mac-crm/pkg/config"
	"dmglab.com/mac-crm/pkg/models"
	"dmglab.com/mac-crm/pkg/service"
	"dmglab.com/mac-crm/pkg/util"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Migration(migration *gorm.DB) {
	if config.GetConfig().DBDriver == "MYSQL" || config.GetConfig().DBDriver == "CLOUDMYSQL" {
		util.CheckAndCreate(migration.Migrator(), &models.App{})
		util.CheckAndCreate(migration.Migrator(), &models.Role{})
		util.CheckAndCreate(migration.Migrator(), &models.Account{})
		util.CheckAndCreate(migration.Migrator(), &models.Customer{})
		util.CheckAndCreate(migration.Migrator(), &models.Attachment{})
		util.CheckAndCreate(migration.Migrator(), &models.CustomField{})

		//meta table
		migration.AutoMigrate(&models.CustomField{})
		migration.AutoMigrate(&models.Customer{})
		migration.AutoMigrate(&models.CustomersMeta{})
		migration.AutoMigrate(&models.Role{})
		migration.AutoMigrate(&models.FieldOption{})

	}
}

func InitalSystemAcc(migration *gorm.DB) {
	if migration.Migrator().HasTable(&models.Role{}) {
		if err := migration.Clauses(clause.OnConflict{DoNothing: true}).Model(&models.Role{}).Create(&models.Role{
			ID:               uuid.UUID{},
			Name:             _const.ACC_SUPER_ADMIN,
			Permissions:      _const.ROLE_SUPER,
			FieldPermissions: "*",
			Levels:           "*",
		}).Error; err != nil {
			service.SysLog.Fatalln(err)
		}
	}
	if migration.Migrator().HasTable(&models.Account{}) {
		if err := migration.Clauses(clause.OnConflict{DoNothing: true}).Model(&models.Account{}).Create(&models.Account{
			ID:          uuid.UUID{},
			DisplayName: _const.ACC_SCHEDULER,
			Username:    "scheduler",
			Status:      "Active",
			IsSystem:    true,
		}).Error; err != nil {
			service.SysLog.Fatalln(err)
		}
	}
	if migration.Migrator().HasTable(&models.App{}) {
		if err := migration.Clauses(clause.OnConflict{DoNothing: true}).Model(&models.App{}).
			Create([]*models.App{
				{
					Setting: _const.ROW_LIMIT,
					Value:   []byte("500"),
					Version: 0,
				},
			}).Error; err != nil {
			service.SysLog.Fatalln(err)
		}
	}
}
func InitalAccount(migration *gorm.DB) {
	var u string
	var p string
	flag.StringVar(&u, "u", "", "admin user")
	flag.StringVar(&p, "p", "", "admin password")
	flag.Parse()
	if len(u) == 0 || len(p) == 0 {
		service.SysLog.Errorln("invalid username and password")
		return
	}
	defaultUUID := uuid.UUID{}.String()
	defaultUUID = defaultUUID[:len(defaultUUID)-1] + "1"
	id, err := uuid.Parse(defaultUUID)
	if err != nil {
		service.SysLog.Errorln(err)
		return
	}
	acc := &models.Account{
		ID:       id,
		Username: u,
		RoleID:   uuid.UUID{},
	}
	acc.SetPassword(p)
	err = migration.Model(&models.Account{}).FirstOrCreate(acc).Error
	if err != nil {
		service.SysLog.Errorln(err)
		return
	}
}
