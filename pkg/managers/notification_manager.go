package managers

import (
	"context"

	"dmglab.com/mac-crm/pkg/config"
	"dmglab.com/mac-crm/pkg/models"
	"dmglab.com/mac-crm/pkg/util"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationQueryParam struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
type NotificationUpdateParam struct {
	IDs    []string `josn:"ids"`
	IsRead bool     `json:"is_read"`
}
type INotificationManager interface {
	Updates(ctx context.Context, param *NotificationUpdateParam) ([]*models.Notification, []*models.NotificationReader, error)
	GetNotifications(ctx context.Context, param *NotificationQueryParam) ([]*models.Notification, []*models.NotificationReader, *util.Pagination, error)
}

type NotificationManager struct {
	config *config.Config
}

func GetNotificationManager() INotificationManager {
	return &NotificationManager{
		config: config.GetConfig(),
	}
}

func (m *NotificationManager) Updates(ctx context.Context, param *NotificationUpdateParam) ([]*models.Notification, []*models.NotificationReader, error) {
	notifications := []*models.Notification{}
	notificationReaders := []*models.NotificationReader{}
	acc, err := util.GetCtxAccount(ctx)
	if err != nil {
		return nil, nil, err
	}
	return notifications, notificationReaders, util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		err := tx.Model(notifications).Where("id in (?) ", param.IDs).Updates(map[string]interface{}{
			"is_read":    param.IsRead,
			"account_id": acc.ID,
		}).Error
		if err != nil {
			return err
		}
		err = tx.Preload("Account").Where("id in (?) ", param.IDs).Find(&notifications).Error
		if err != nil {
			return err
		}
		return tx.Preload("Reader").Where("reader_id = ? ", acc.ID).Where("notification_id in (?) ", param.IDs).Find(&notificationReaders).Error
	})
}
func (m *NotificationManager) GetNotifications(ctx context.Context, param *NotificationQueryParam) ([]*models.Notification, []*models.NotificationReader, *util.Pagination, error) {
	pagin := &util.Pagination{
		Limit: param.Limit,
		Page:  param.Page,
	}
	notifications := []*models.Notification{}
	notificationReaders := []*models.NotificationReader{}
	acc, err := util.GetCtxAccount(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	return notifications, notificationReaders, pagin, util.GetCtxTx(ctx, func(tx *gorm.DB) error {
		err := tx.Preload("Account").Order("created_at desc").Scopes(util.PaginationScope(notifications, pagin, tx)).Find(&notifications).Error
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		notiIDs := make([]uuid.UUID, len(notifications))
		for i, noti := range notifications {
			notiIDs[i] = noti.ID
		}
		err = tx.Preload("Reader").Where("reader_id = ? ", acc.ID).Where("notification_id in (?)", notiIDs).Find(&notificationReaders).Error
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	})
}
