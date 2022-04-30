package models

import (
	"encoding/json"
	"time"

	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID        uuid.UUID      `gorm:"type:varchar(36);primaryKey;" json:"id"`
	SKU       string         `gorm:"type:varchar(36);" `
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Event     ProductEvent   `gorm:"<-:false;foreignKey:ProductID;references:ID" json:"-"`
}

func (product *Product) ToJSON() ([]byte, error) {

	return json.Marshal(product)
}

func (product *Product) BeforeCreate(tx *gorm.DB) (err error) {
	product.ID = uuid.New()
	return
}

func (product *Product) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, product)
	return
}
func (product *Product) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, product)
	return
}
