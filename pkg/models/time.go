package models

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Date struct {
	time.Time
}

func (date *Date) Format(layout string) string {
	if date == nil {
		return ""
	}
	return date.Time.Format(layout)
}

func (date Date) GormDataType() string {
	return "date"
}

func (date *Date) Scan(v interface{}) error {
	date.Time = v.(time.Time)
	return nil
}
func (date *Date) UnmarshalJSON(bytes []byte) error {

	var raw int64
	err := json.Unmarshal(bytes, &raw)

	if err != nil {
		fmt.Printf("error decoding timestamp: %s\n", err)
		return err
	}
	if raw > 9999999999 {
		raw = raw / 1000
	}
	date.Time = time.Unix(raw, 0)
	return nil
}

func (d Date) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "CONVERT(DATE,?)",
		Vars: []interface{}{d.Time.Format("2006-01-02")},
	}
}

func (date *Date) ToNullableTimeStamp() *int64 {
	if date == nil {
		return nil
	}
	if ts := date.Unix(); ts > 0 {
		return &ts
	}
	return nil
}

type Timestamp struct {
	*time.Time
}

func (t *Timestamp) UnmarshalJSON(bytes []byte) error {
	var raw int64
	err := json.Unmarshal(bytes, &raw)

	if err != nil {
		fmt.Printf("error decoding timestamp: %s\n", err)
		return err
	}

	tt := time.Unix(raw, 0)
	t.Time = &tt
	return nil
}

func (t Timestamp) GormDataType() string {
	return "timestamp NULL"
}

func (t *Timestamp) Scan(v interface{}) error {

	time, ok := v.(time.Time)
	if !ok {
		return errors.New("invalid time")
	}
	t.Time = &time
	return nil
}
func (t Timestamp) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t, nil

}
func (t Timestamp) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "?",
		Vars: []interface{}{t.Time.Unix()},
	}
}

func (t *Timestamp) ToInt64() *int64 {
	if t == nil {
		return nil
	}
	if ts := t.Unix(); ts > 0 {
		return &ts
	}
	return nil
}
