package models

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

//Token
type NullableUUID struct {
	UUID  uuid.UUID
	Valid bool // Valid is true if String is not NULL
}

func (n *NullableUUID) String() string {
	if !n.Valid {
		return ""
	}
	return n.UUID.String()
}
func (n *NullableUUID) Scan(value interface{}) error {
	if value == nil {
		n.UUID, n.Valid = uuid.UUID{}, false
		return nil
	}
	n.Valid = true
	v := value.(string)
	err := validator.New().Var(v, "uuid")
	if err != nil {
		return err
	}

	n.UUID, err = uuid.Parse(v)
	return err
}

func (n NullableUUID) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.String, nil
}

func (n NullableUUID) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n)
	}
	return json.Marshal(nil)
}

func (n *NullableUUID) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.UUID)
	if err == nil {
		n.Valid = true
	}
	return err
}

//Token
type Token sql.NullString

// Scan implements the Scanner interface.
func (n *Token) Scan(value interface{}) error {
	return (*sql.NullString)(n).Scan(value)
}

func (n Token) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.String, nil
}

func (n Token) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	}
	return json.Marshal(nil)
}

func (n *Token) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.String)
	if err == nil {
		n.Valid = true
	}
	return err
}

//JSON
type JSON string

func (j *JSON) Scan(value interface{}) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	result := json.RawMessage{}
	err = json.Unmarshal(bytes, &result)
	*j = JSON(result)
	return err
}

// Value return json value, implement driver.Valuer interface
func (j JSON) Value() ([]byte, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.RawMessage(j).MarshalJSON()
}

type NullableTime sql.NullTime

// Scan implements the Scanner interface.
func (n *NullableTime) Scan(value interface{}) error {
	return (*sql.NullTime)(n).Scan(value)
}

// Value implements the driver Valuer interface.
func (n NullableTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

func (n NullableTime) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.Time)
	}
	return json.Marshal(nil)
}

func (n *NullableTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		n.Valid = false
		return nil
	}
	err := json.Unmarshal(b, &n.Time)
	if err == nil {
		n.Valid = true
	}
	return err
}
