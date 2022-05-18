package entities

import (
	"dmglab.com/mac-crm/pkg/models"
)

type FieldOption struct {
	ID      string         `json:"id"`
	Name    *MultiLangText `json:"name"`
	FieldID string         `json:"field_id"`
}

func NewFieldOptinoEntity(option *models.FieldOption) *FieldOption {
	if option == nil {
		return &FieldOption{}
	}
	return &FieldOption{
		ID:      option.ID.String(),
		Name:    NewMultiLangTextEntity(option.Name),
		FieldID: option.FieldID.String(),
	}
}

func NewFieldOptinoListEntity(options []*models.FieldOption) []*FieldOption {
	optionList := make([]*FieldOption, len(options))
	for i, option := range options {
		optionList[i] = NewFieldOptinoEntity(option)
	}
	return optionList
}
