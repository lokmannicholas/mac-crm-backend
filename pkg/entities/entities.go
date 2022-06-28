package entities

import (
	"time"
)

func timeToMilli(target *time.Time) *int64 {
	if target == nil {
		return nil
	}
	timeStamp := target.UnixMilli()
	return &timeStamp
}

type Entity interface{}

type List struct {
	Columns []string    `json:"columns"`
	Total   int64       `json:"total"`
	Data    interface{} `json:"data"`
}
