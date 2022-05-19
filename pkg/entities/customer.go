package entities

import (
	"bytes"
	"context"
	"encoding/gob"
	"reflect"
	"strings"

	"dmglab.com/mac-crm/pkg/managers"
	"dmglab.com/mac-crm/pkg/models"
)

type Customer struct {
	ID        string `json:"id"`
	Code      string `json:"code"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	OtherName string `json:"other_name"`
	Phone     string `json:"phone"`
	IDNo      string `json:"id_no"`
	Title     string `json:"title"`
	Address   string `json:"address"`
	Status    string `json:"status" enums:"Active,Disable" default:"Active"`
	Remarks   string `json:"remarks"`
	Meta      []Meta `json:"meta"`
}

func (b Customer) Fields() []string {
	val := reflect.ValueOf(b)
	noOfFields := val.Type().NumField()
	fields := make([]string, noOfFields)
	for i := 0; i < noOfFields; i++ {
		fields[i] = val.Type().Field(i).Tag.Get("json")
	}
	return fields
}
func NewCustomerEntity(customer *models.Customer, ctx context.Context) *Customer {
	if customer == nil {
		return &Customer{}
	}

	idNo := "******"
	if len(customer.GetIDNo()) > 4 {
		idNo = customer.GetIDNo()[0:4] + "******"
	}
	metaArray := make([]Meta, len(customer.Meta))
	for i, meta := range customer.Meta {
		buf := bytes.NewBuffer(meta.Val)
		dec := gob.NewDecoder(buf)
		v := ""
		ent := *NewMetaEntity(meta.Meta)
		ent.Val = ""
		if err := dec.Decode(&v); err == nil {
			if ent.DataType == "multiple" {
				optionIds := strings.Split(v, ";")
				options, err := managers.FindByIds(ctx, optionIds)
				if err == nil {
					ent.Val = NewFieldOptinoListEntity(options)
				}
			} else {
				ent.Val = v
			}
		}
		metaArray[i] = ent
	}
	return &Customer{
		ID:        customer.ID.String(),
		Code:      customer.Code,
		FirstName: customer.FirstName,
		LastName:  customer.LastName,
		OtherName: customer.OtherName,
		Phone:     customer.Phone,
		IDNo:      idNo,
		Title:     customer.Title,
		Address:   customer.Adderess,
		Status:    customer.Status,
		Remarks:   customer.Remarks,
		Meta:      metaArray,
	}
}

func NewCustomerListEntity(total int64, customers []*models.Customer, ctx context.Context) *List {
	customerList := make([]*Customer, len(customers))
	for i, customer := range customers {
		customerList[i] = NewCustomerEntity(customer, ctx)
	}
	return &List{
		Columns: Customer{}.Fields(),
		Total:   total,
		Data:    customerList,
	}

}
