package pubsub

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"dmglab.com/mac-crm/pkg/util"
	"github.com/google/uuid"
)

type Command struct {
	Token   string
	Command string
	Data    string
}

type Message struct {
	Topic     string     `json:"topic"`
	Data      []rune     `json:"data"`
	Timestamp int64      `json:"timestamp"`
	UserID    *uuid.UUID `json:"user_id"`
}

type PubSubService struct {
	messagesQueue chan Message
	subscriber    map[string][]string
	mu            sync.RWMutex
}

var service = &PubSubService{
	messagesQueue: make(chan Message, 3000),
	subscriber:    make(map[string][]string),
}

func GetPubSubService() *PubSubService {
	return service
}

func (ser *PubSubService) Register(topic string, subscriber string) {
	if val, ok := ser.subscriber[topic]; ok {
		ser.subscriber[topic] = append(val, subscriber)
	} else {
		ser.subscriber[topic] = []string{subscriber}
	}
}

func (ser *PubSubService) Push(ctx context.Context, topic string, msg interface{}) {
	ser.mu.RLock()
	defer ser.mu.RUnlock()
	// var data Message

	b, err := json.Marshal(msg)
	if err != nil {
		return
	}
	// data = Message{
	// 	Topic: topic,
	// 	Data:  []rune(string(b)),
	// }
	var userID *uuid.UUID
	if ctx != nil {
		acc, _ := util.GetCtxAccount(ctx)
		if acc != nil {
			userID = &acc.ID
		}
	}
	ser.messagesQueue <- Message{
		Topic:     topic,
		Data:      []rune(string(b)),
		UserID:    userID,
		Timestamp: time.Now().Unix(),
	}
	// if ch, ok := ser.messagesQueue; ok {

	// 	ch <- data
	// } else {
	// 	ser.messagesQueue[topic] = make(chan Message, 3000)
	// 	ser.messagesQueue[topic] <- data
	// }

}

func (ser *PubSubService) Subscribe(f func(msg *Message) error) {
	for {
		ser.mu.RLock()
		select {
		case m := <-ser.messagesQueue:
			{
				f(&m)
				ser.mu.RUnlock()
			}
		default:
			ser.mu.RUnlock()
		}
	}
}
