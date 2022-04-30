package entities

import (
	"dmglab.com/mac-crm/pkg/models"
)

type CustomField struct {
	ID           string         `json:"id"`
	UniqueKey    string         `json:"unique_key"`
	CustomObject string         `json:"custom_object"`
	FieldName    *MultiLangText `json:"field_name"`
	FieldType    string         `json:"field_type"`
	Status       string         `json:"status"`
	Remarks      string         `json:"remarks"`
}

func NewCustomFieldEntity(cus *models.CustomField) *CustomField {
	if cus == nil {
		return &CustomField{}
	}

	return &CustomField{
		ID:           cus.ID.String(),
		CustomObject: cus.CustomObject,
		FieldName:    NewMultiLangTextEntity(cus.FieldName),
		FieldType:    cus.FieldType,
		UniqueKey:    cus.UniqueKey,
		Status:       cus.Status,
		Remarks:      cus.Remarks,
	}
}

func NewCustomFieldListEntity(total int64, cuss []*models.CustomField) *List {
	cusList := make([]*CustomField, len(cuss))
	for i, cus := range cuss {
		cusList[i] = NewCustomFieldEntity(cus)
	}
	return &List{
		Total: total,
		Data:  cusList,
	}

}
