package models

import (
	"encoding/json"
	"time"

	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Storage struct {
	ID            uuid.UUID   `gorm:"type:varchar(36);primaryKey;" json:"id"`
	DisplayID     string      `gorm:"<-:create;type:varchar(36);" json:"display_id"`
	Code          string      `gorm:"column:code" json:"code"`
	Categories    []*Category `gorm:"<-:false;many2many:storage_category;foreignKey:ID;" json:"categories"`
	Features      []*Feature  `gorm:"<-:false;many2many:storage_feature;foreignKey:ID;" json:"features" `
	RentStartDate *time.Time  `json:"rent_start_date,omitempty"`
	RentEndDate   *time.Time  `json:"rent_end_date,omitempty"`
	RentType      string      `json:"rent_type"`
	Status        string      `json:"status"` //setAvailable, Disable, Rent, Pending, Expired
	//belongs to
	Customer     *Customer  `gorm:"<-:false;foreignKey:CustomerID" json:"-"`
	CustomerID   *uuid.UUID `gorm:"type:varchar(36);" json:"customer_id,omitempty"`
	MonthlyPrice *float32   `json:"monthly_price,omitempty"`
	YearlyPrice  *float32   `json:"yearly_price,omitempty"`
	DailyPrice   *float32   `json:"daily_price,omitempty"`
	//has many
	StorageType string        `json:"storage_type" ` //PrivateStorage/ InOutStorage
	RentRecords []*RentRecord `gorm:"many2many:storage_rent_records;foreignKey:ID;joinForeignKey:StorageID;References:ID;JoinReferences:RentRecordID" json:"-"`

	ContractID *uuid.UUID `gorm:"type:varchar(36);" json:"contract_id,omitempty"`
	Contract   *Contract  `gorm:"<-:false;foreignKey:ContractID" json:"-"`

	OrderID      *uuid.UUID      `gorm:"type:varchar(36);" json:"order_id,omitempty"`
	CurrentOrder *RentalOrder    `gorm:"<-:false;foreignKey:OrderID" json:"-"`
	Events       []*StorageEvent `gorm:"<-:create;foreignKey:StorageID;references:ID" json:"-"`

	Remarks    string       `json:"remarks,omitempty"`
	ScheduleID *uuid.UUID   `gorm:"<-:update;type:varchar(36);" json:"-"`
	Schedule   *ScheduleJob `gorm:"foreignKey:ScheduleID" json:"-"`

	BranchID  uuid.UUID      `gorm:"<-:create;type:varchar(36);index:idx_branch_id" json:"branch_id"`
	Branch    *Branch        `gorm:"foreignKey:BranchID" json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Event     StorageEvent   `gorm:"<-:false;foreignKey:StorageID;references:ID" json:"-"`
}

func (storage *Storage) ToJSON() ([]byte, error) {

	return json.Marshal(storage)
}

func (storage *Storage) BeforeCreate(tx *gorm.DB) (err error) {
	storage.ID = uuid.New()
	storage.Status = _const.STORAGE_AVAILABLE
	storage.DisplayID = storage.Code
	return
}

func (storage *Storage) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, storage)
	return
}
func (storage *Storage) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, storage)
	return
}

func (storage *Storage) Copy(newStorage *Storage) {
	newStorage.RentStartDate = storage.RentStartDate
	newStorage.RentEndDate = storage.RentEndDate
	newStorage.RentType = storage.RentType
	newStorage.Status = storage.Status
	newStorage.Customer = storage.Customer
	newStorage.CustomerID = storage.CustomerID
	newStorage.OrderID = storage.OrderID
	newStorage.CurrentOrder = storage.CurrentOrder
}
func (storage *Storage) ResumeAvailable() {
	storage.RentStartDate = nil
	storage.RentEndDate = nil
	storage.RentType = ""
	storage.Status = _const.STORAGE_AVAILABLE

	storage.Customer = nil
	storage.CustomerID = nil
	storage.OrderID = nil
	storage.CurrentOrder = nil
}
