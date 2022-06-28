package entities

import (
	"time"
)

func timeToMilli(target *time.Time) *int64 {
	if target == nil {
		return nil
	}
	timeStamp := target.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
	return &timeStamp
}

type Entity interface{}

type List struct {
	Columns []string    `json:"columns"`
	Total   int64       `json:"total"`
	Data    interface{} `json:"data"`
}
