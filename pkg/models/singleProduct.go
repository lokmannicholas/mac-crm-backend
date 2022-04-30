package models

import (
	"encoding/json"
	"time"

	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SingleProduct struct {
	ID          uuid.UUID          `gorm:"type:varchar(36);primaryKey;" json:"id"`
	SKU         string             `gorm:"type:varchar(36);unique" json:"sku"`
	ProductType string             `gorm:"type:varchar(12);" json:"product_type"`
	Name        *MultiLangText     `gorm:"embedded;embeddedPrefix:name_" json:"name"`
	Customer    *Customer          `gorm:"<-:false;foreignKey:CustomerID" json:"-"`
	CustomerID  *uuid.UUID         `gorm:"type:varchar(36);" json:"customer_id,omitempty"`
	Price       float32            `json:"price"`
	Status      string             `json:"status"`
	Remarks     string             `json:"remarks"`
	CreatedAt   time.Time          `json:"-"`
	UpdatedAt   time.Time          `json:"-"`
	DeletedAt   gorm.DeletedAt     `gorm:"index" json:"-"`
	Event       SingleProductEvent `gorm:"<-:false;foreignKey:SingleProductID;references:ID" json:"-"`
}

func (product *SingleProduct) ToJSON() ([]byte, error) {
	return json.Marshal(product)
}

func (product *SingleProduct) BeforeCreate(tx *gorm.DB) (err error) {
	product.ID = uuid.New()
	return
}

func (product *SingleProduct) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, product)
	return
}
func (product *SingleProduct) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, product)
	return
}
