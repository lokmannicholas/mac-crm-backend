package notificationcenter

import (
	"dmglab.com/mac-crm/pkg/collections"
	"dmglab.com/mac-crm/pkg/models"
	"dmglab.com/mac-crm/pkg/service"
	"dmglab.com/mac-crm/pkg/service/pubsub"
	"gorm.io/gorm"
)

type INotificationCenter interface {
	SetNotice()
}

type NotificationCenter struct {
	DB *gorm.DB
}

var noti *NotificationCenter

func GetNotificationCenter() INotificationCenter {
	if noti == nil {
		noti = &NotificationCenter{
			DB: collections.GetCollection().DB,
		}
	}

	return noti
}

func (noti *NotificationCenter) SetNotice() {
	ser := pubsub.GetPubSubService()
	go ser.Subscribe(func(msg *pubsub.Message) error {
		if msg != nil {
			notification := new(models.Notification)
			if msg.UserID != nil {
				notification.AccountID = *msg.UserID
			}
			notification.EventType = msg.Topic
			notification.Content = []byte(string(msg.Data))
			err := noti.DB.Transaction(func(tx *gorm.DB) error {
				return tx.Create(notification).Error
			})
			if err != nil {
				service.SysLog.Errorln(err)
			}
		}
		return nil
	})
}
