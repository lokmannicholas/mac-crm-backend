package entities

import (
	"dmglab.com/mac-crm/pkg/models"
)

type FieldOption struct {
	Name *MultiLangText `json:"name"`
}

func NewFieldOptinoEntity(option *models.FieldOption) *FieldOption {
	if option == nil {
		return &FieldOption{}
	}
	return &FieldOption{
		Name: NewMultiLangTextEntity(option.Name),
	}
}

func NewFieldOptinoListEntity(options []*models.FieldOption) []*FieldOption {
	optionList := make([]*FieldOption, len(options))
	for i, option := range options {
		optionList[i] = NewFieldOptinoEntity(option)
	}
	return optionList
}
