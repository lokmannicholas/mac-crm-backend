package models

import (
	"time"

	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"gorm.io/gorm"
)

type App struct {
	Setting   string `gorm:"type:varchar(36);primaryKey;"`
	Value     []byte
	Version   int       `gorm:"primaryKey;"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time `gorm:"<-:update"`
}

func (App) TableName() string {
	return "app_config"
}

func (app *App) BeforeCreate(tx *gorm.DB) (err error) {
	app.CreatedAt = time.Now()
	return
}

func (app *App) BeforeSave(tx *gorm.DB) (err error) {
	app.UpdatedAt = time.Now()
	return
}

func (app *App) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, app)
	return
}
