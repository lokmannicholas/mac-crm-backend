package pubsub

import (
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	ser := GetPubSubService()
	// ser.Register("test1", "")
	t.Log("test1")
	go ser.Subscribe(func(msg *Message) error {
		if msg != nil {
			t.Logf("%s:%s", msg.Topic, string(msg.Data))
		}
		return nil
	})
	for {
		ser.Push(nil, "test1", map[string]interface{}{
			"time": "now",
		})
		time.Sleep(time.Second * 1)
	}
}
