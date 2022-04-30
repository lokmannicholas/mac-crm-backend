package models

import (
	"encoding/json"

	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"dmglab.com/mac-crm/pkg/util/encrypt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	ID          uuid.UUID    `gorm:"type:varchar(36);primaryKey;" json:"id"`
	DisplayName string       `json:"display_name"`
	Username    string       `gorm:"unique" json:"user_name"`
	Password    string       `json:"-"`
	RoleID      uuid.UUID    `gorm:"type:varchar(36);" json:"role_id"`
	Role        Role         `gorm:"<-:false;foreignKey:RoleID" json:"role"`
	Status      string       `json:"status"`
	IsSystem    bool         `gorm:"default:false;" json:"is_system"`
	LastLogin   NullableTime `json:"-"`
}

func (acc *Account) IsScheduler() bool {
	return acc.DisplayName == _const.ACC_SCHEDULER
}
func (acc *Account) BeforeCreate(tx *gorm.DB) (err error) {
	if (acc.ID == uuid.UUID{}) {
		acc.ID = uuid.New()
	}
	acc.SetPassword(acc.Password)
	acc.SetActive()
	return
}

func (acc *Account) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, acc)
	return
}

func (acc *Account) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, acc)
	return
}

func (acc *Account) SetDisable() *Account {
	acc.Status = "Disable"
	return acc
}

func (acc *Account) SetPassword(password string) string {
	acc.Password = encrypt.MD5Hash(password)
	return acc.Password
}
func (acc *Account) SetActive() *Account {
	acc.Status = "Active"
	return acc
}
func (acc *Account) ToJSON() string {
	b, err := json.Marshal(acc)
	if err != nil {
		return ""
	}
	return string(b)
}
