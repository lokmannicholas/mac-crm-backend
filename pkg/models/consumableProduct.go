package models

import (
	"encoding/json"
	"time"

	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ConsumableProduct struct {
	ID          uuid.UUID              `gorm:"type:varchar(36);primaryKey;" json:"id"`
	SKU         string                 `gorm:"type:varchar(36);unique" json:"sku"`
	ProductType string                 `gorm:"type:varchar(12);" json:"product_type"`
	Name        *MultiLangText         `gorm:"embedded;embeddedPrefix:name_" json:"name"`
	Price       float32                `json:"price"`
	Inventory   int                    `json:"inventory"`
	Status      string                 `json:"status"`
	Remarks     string                 `json:"remarks"`
	CreatedAt   time.Time              `json:"-"`
	UpdatedAt   time.Time              `json:"-"`
	DeletedAt   gorm.DeletedAt         `gorm:"index" json:"-"`
	Event       ConsumableProductEvent `gorm:"<-:false;foreignKey:ConsumableProductID;references:ID" json:"-"`
}

func (product *ConsumableProduct) ToJSON() ([]byte, error) {
	return json.Marshal(product)
}

func (product *ConsumableProduct) BeforeCreate(tx *gorm.DB) (err error) {
	product.ID = uuid.New()
	return
}

func (product *ConsumableProduct) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, product)
	return
}
func (product *ConsumableProduct) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, product)
	return
}
