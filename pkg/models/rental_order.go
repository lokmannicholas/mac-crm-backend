package models

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RentalOrder struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:varchar(36);primaryKey;" json:"id"`
	OrderNo   string    `gorm:"<-:create;unique;type:varchar(36)" json:"order_no"` //      `gorm:"<-:create;unique;type:MEDIUMINT UNSIGNED NOT NULL AUTO_INCREMENT"`
	OrderType string    `json:"order_type"`
	OrderDate time.Time `json:"order_date"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`

	ActualStartDate time.Time `json:"actual_start_date"`
	ActualEndDate   time.Time `json:"actual_end_date"`

	RentType   string `json:"rent_type"` //month year day
	RentPeriod int    `json:"rent_period"`
	Status     string `json:"status"`

	OriginalPrice   float32            `gorm:"<-:create;" json:"original_price"`
	DiscountAmount  float32            `gorm:"<-:create;" json:"discount_amount"`
	FinalPrice      float32            `json:"final_price"`
	IsExpired       bool               `json:"is_expired"`
	TerminationDate *time.Time         `json:"termination_date"`
	TransferDate    *time.Time         `json:"transfer_date"`
	TransferAmount  *float32           `json:"transfer_amount"`
	Remarks         string             `json:"remarks"`
	Item            []*RentalOrderItem `gorm:"<-:create;foreignKey:OrderID;references:ID" json:"item"`
	//belongs to
	Customer   Customer  `gorm:"<-:false;foreignKey:CustomerID" json:"-"`
	CustomerID uuid.UUID `gorm:"type:varchar(36);" json:"customer_id"`

	//belongs to
	Staff     Account   `gorm:"<-:false;foreignKey:StaffID" json:"-"`
	StaffID   uuid.UUID `gorm:"type:varchar(36);" json:"staff_id"`
	StaffName string    `json:"staff_name"`

	RefFrom *uuid.UUID `gorm:"type:varchar(36);" json:"ref_from"`
	RefType *string    `gorm:"type:varchar(16);" json:"ref_type"`

	Invoices  []*Invoice `gorm:"foreignKey:OrderID;references:ID" json:"invoices"`
	AutoRenew bool       `json:"auto_renew"`

	BranchID uuid.UUID `gorm:"type:varchar(36);" json:"branch_id"`
	Branch   Branch    `gorm:"<-:false;foreignKey:BranchID" json:"-"`

	ContractID uuid.UUID `gorm:"type:varchar(36);index" json:"contract_id"`
	Contract   Contract  `gorm:"<-:false;foreignKey:ContractID" json:"-"`

	Event RentalOrderEvent `gorm:"<-:false;foreignKey:OrderID;references:ID" json:"-"`
}

func (order *RentalOrder) BeforeCreate(tx *gorm.DB) (err error) {
	order.ID = uuid.New()
	// order.StartDate = time.Date(order.StartDate.Year(), order.StartDate.Month(), order.StartDate.Day(), 0, 0, 0, 0, order.StartDate.Location())
	// order.EndDate = time.Date(order.EndDate.Year(), order.EndDate.Month(), order.EndDate.Day(), 0, 0, 0, 0, order.EndDate.Location())

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	randomNumber := r1.Intn(100-1) + 1
	order.OrderNo = fmt.Sprintf("R%s%d", order.OrderDate.Format("20060102150405"), randomNumber)
	order.AutoRenew = true

	return
}

func (order *RentalOrder) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, order)
	return
}

func (order *RentalOrder) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, order)
	return
}

func (order *RentalOrder) ToJSON() ([]byte, error) {

	return json.Marshal(order)
}

//mssql
type RentalOrderItem struct {
	OrderID        uuid.UUID `gorm:"type:varchar(36);primaryKey;" json:"order_id"`
	ItemType       string    `json:"item_type"`
	ItemID         uuid.UUID `gorm:"type:varchar(36);primaryKey;" json:"item_id"`
	OriginalPrice  float32   `gorm:"<-:create;" json:"original_price"`
	DiscountAmount float32   `gorm:"<-:create;" json:"discount_amount"`
	FinalPrice     float32   `json:"final_price"`
	Storage        *Storage  `gorm:"-"`
	Item           []byte
}

func (o *RentalOrderItem) BeforeCreate(tx *gorm.DB) (err error) {
	if len(o.Item) == 0 {
		b, err := o.Storage.ToJSON()
		if err != nil {
			return err
		}
		o.Item = b
	}

	return
}

type RentalOrderEvent struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:varchar(36);primaryKey;"`
	Event       string
	OrderID     uuid.UUID `gorm:"type:varchar(36);"`
	Update      []byte    `gorm:"type:text;"`
	TriggerByID uuid.UUID `gorm:"type:varchar(36);"`
	TriggerBy   Account   `gorm:"<-:false;foreignKey:TriggerByID"`
}

func (orderEvent *RentalOrderEvent) BeforeCreate(tx *gorm.DB) (err error) {
	orderEvent.ID = uuid.New()
	return
}

func (orderEvent *RentalOrderEvent) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, orderEvent)
	return
}

func (orderEvent *RentalOrderEvent) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, orderEvent)
	return
}
