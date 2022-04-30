package models

import (
	"time"

	"dmglab.com/mac-crm/pkg/service"
	_const "dmglab.com/mac-crm/pkg/util/const"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScheduleJob struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:varchar(36);primaryKey;"`
	ScheduleType string
	ScheduleDate time.Time
	BranchID     string
	StorageID    string
	Status       string
}
type Schedule struct {
	ID           uuid.UUID `gorm:"type:varchar(36);primaryKey;"`
	ScheduleType string
	ScheduleDate time.Time
	StorageID    string
	BranchID     string
	StorageCode  string
	Status       string
}

func (scheduleJob *ScheduleJob) BeforeCreate(tx *gorm.DB) (err error) {
	scheduleJob.ID = uuid.New()
	scheduleJob.Status = _const.SCHEDULE_ACTIVE
	//scheduleJob.ScheduleType = _const.SCHEDULE_TYPE_DAILY
	return
}

func (scheduleJob *ScheduleJob) AfterCreate(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_CREATE, scheduleJob)
	return
}

func (scheduleJob *ScheduleJob) AfterSave(tx *gorm.DB) (err error) {
	service.GetAuditLogger().InfoLog(tx.Statement.Context, _const.AUDIT_UPDATE, scheduleJob)
	return
}
