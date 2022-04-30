package entities

import (
	"encoding/json"

	"dmglab.com/mac-crm/pkg/models"
)

type StorageEvent struct {
	ID        string                 `json:"id"`
	Event     string                 `json:"event"`
	StorageID string                 `json:"storage_id"`
	Update    map[string]interface{} `json:"update"`
	TriggerBy *Account               `json:"trigger_by"`
	CreatedAt int64                  `json:"created_at"`
}

func NewStorageEventEntity(event *models.StorageEvent) *StorageEvent {
	m := map[string]interface{}{}
	if event.Update != nil && len(event.Update) > 0 {
		json.Unmarshal(event.Update, &m)
	}
	return &StorageEvent{
		ID:        event.ID.String(),
		Event:     event.Event,
		StorageID: event.StorageID.String(),
		Update:    m,
		TriggerBy: NewAccountEntity(&event.TriggerBy),
		CreatedAt: event.CreatedAt.Unix(),
	}
}
func NewStorageEventListEntity(total int64, events []*models.StorageEvent) *List {
	storageEventList := make([]*StorageEvent, len(events))
	for i, event := range events {
		storageEventList[i] = NewStorageEventEntity(event)
	}
	return &List{
		Total: total,
		Data:  storageEventList,
	}

}

type OrderEvent struct {
	ID        string                 `json:"id"`
	Event     string                 `json:"event"`
	OrderID   string                 `json:"order_id"`
	Update    map[string]interface{} `json:"update"`
	TriggerBy *Account               `json:"trigger_by"`
	CreatedAt int64                  `json:"created_at"`
}

func NewOrderEventEntity(event *models.RentalOrderEvent) *OrderEvent {
	m := map[string]interface{}{}
	if event.Update != nil && len(event.Update) > 0 {
		json.Unmarshal(event.Update, &m)
	}
	return &OrderEvent{
		ID:        event.ID.String(),
		Event:     event.Event,
		OrderID:   event.OrderID.String(),
		Update:    m,
		TriggerBy: NewAccountEntity(&event.TriggerBy),
		CreatedAt: event.CreatedAt.Unix(),
	}
}
func NewOrderEventListEntity(total int64, events []*models.RentalOrderEvent) *List {
	orderEventList := make([]*OrderEvent, len(events))
	for i, event := range events {
		orderEventList[i] = NewOrderEventEntity(event)
	}
	return &List{
		Total: total,
		Data:  orderEventList,
	}

}
