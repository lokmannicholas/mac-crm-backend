package models

type Meta struct {
	Key      string `gorm:"type:varchar(18);primaryKey;" json:"key"`
	DataType string `json:"data_type"`
	Val      string `json:"val"`
}
