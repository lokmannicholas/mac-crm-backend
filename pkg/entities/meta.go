package entities

import (
	"dmglab.com/mac-crm/pkg/models"
)

type Meta struct {
	Key      string      `json:"key"`
	DataType string      `json:"data_type"`
	Val      interface{} `json:"val"`
}

func NewMetaEntity(meta *models.Meta) *Meta {

	return &Meta{
		Key:      meta.Key,
		DataType: meta.DataType,
		Val:      meta.Val,
	}
}
