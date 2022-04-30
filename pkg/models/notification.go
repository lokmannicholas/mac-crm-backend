package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	ID        uuid.UUID `gorm:"type:varchar(36);primaryKey;"`
	AccountID uuid.UUID `gorm:"type:varchar(36);"`
	Account   *Account  `gorm:"<-:false;foreignKey:AccountID"`
	Content   []byte    `gorm:"type:BLOB;" json:"content"`
	CreatedAt time.Time
	EventType string
}
type NotificationReader struct {
	NotificationID uuid.UUID `gorm:"type:varchar(36);primaryKey;"`
	ReaderID       uuid.UUID `gorm:"type:varchar(36);primaryKey"`
	Reader         *Account  `gorm:"<-:false;foreignKey:ReaderID"`
	IsRead         bool
	ReadAt         time.Time
}

func (notice *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	notice.ID = uuid.New()
	notice.CreatedAt = time.Now()
	return
}
