package swagger

import "dmglab.com/mac-crm/pkg/entities"

type CustomField struct {
	CustomField *entities.CustomField `json:"custom_field"`
}

type CustomFields struct {
	CustomFields *CustomFieldList `json:"custom_fields"`
}

type CustomFieldList struct {
	List
	Data []*entities.CustomField `json:"data"`
}
