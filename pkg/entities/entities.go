package entities

import (
	"time"
)

func TimeToInt64(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	var i *int64
	timeStamp := t.Unix()
	i = &timeStamp
	return i
}

type Entity interface{}

type List struct {
	Columns []string    `json:"columns"`
	Total   int64       `json:"total"`
	Data    interface{} `json:"data"`
}
