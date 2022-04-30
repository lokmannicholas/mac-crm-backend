package entities

import "dmglab.com/mac-crm/pkg/models"

type MultiLangText struct {
	En string `json:"en"`
	Zh string `json:"zh"`
	Ch string `json:"ch"`
}

func NewMultiLangTextEntity(mul *models.MultiLangText) *MultiLangText {
	if mul == nil {
		return &MultiLangText{}
	}
	return &MultiLangText{
		En: mul.En,
		Zh: mul.Zh,
		Ch: mul.Ch,
	}
}
