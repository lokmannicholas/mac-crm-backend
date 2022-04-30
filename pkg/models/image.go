package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Image struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:varchar(36);primaryKey;"`
	ImageType string
	Data      []byte    `gorm:"type:BLOB;"`
	RefID     uuid.UUID `gorm:"type:varchar(36);"`
}
