package eventLogger

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type EventLogger struct {
	client       *mongo.Client
	messageQueue chan []byte
	end          chan bool
}

var eventLogger *EventLogger

func GetInstance() *EventLogger {
	if eventLogger == nil {
		eventLogger = &EventLogger{
			messageQueue: make(chan []byte, 100),
		}
	}
	return eventLogger
}

func connect() *mongo.Client {
	ctx := context.Background()
	// eventLoggerConfig := config.GetConfig().EventLogger
	uri := fmt.Sprintf("mongodb://%s:%d", "", 0) // eventLoggerConfig.Host,
	// eventLoggerConfig.Port,

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	return client

}
func (e *EventLogger) Stop() {
	e.end <- true

}
func (e *EventLogger) sender(msg []byte) {
	e.messageQueue <- msg

}
func (e *EventLogger) receiver() {
	for {
		select {
		case msg := <-e.messageQueue:
			fmt.Println(string(msg))
		case <-e.end:
			return
		default:

		}
	}

}

func (e *EventLogger) Listen() {
	defer func() {
		if r := recover(); r != nil {
			e.end <- true
		}
	}()
	e.client = connect()
	go e.receiver()

}
