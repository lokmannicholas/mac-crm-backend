package models

import "encoding/json"

type MultiLangText struct {
	En string `json:"en"`
	Zh string `json:"zh"`
	Ch string `json:"ch"`
}

func (mul *MultiLangText) ToJSON() []byte {
	b, err := json.Marshal(mul)
	if err != nil {
		return nil
	}
	return b
}
