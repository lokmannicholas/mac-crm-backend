package ws

import (
	"encoding/json"
	"net/http"
	"sync"

	"dmglab.com/mac-crm/pkg/service"
	"dmglab.com/mac-crm/pkg/service/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Connection struct {
	Socket *websocket.Conn
	mu     sync.Mutex
}

func (c *Connection) Send(message interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.Socket.WriteJSON(message)
}

type IWebSocketController interface {
	Notification(c *gin.Context)
}
type WebSocketController struct {
	PubSub  *pubsub.PubSubService
	Sockets map[string][]*Connection
	Topics  chan string
}

func NewWebSocketController() IWebSocketController {
	return &WebSocketController{
		PubSub:  pubsub.GetPubSubService(),
		Sockets: make(map[string][]*Connection),
		Topics:  make(chan string, 100),
	}
}

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Response struct {
	Topic     string `json:"topic"`
	Data      string `json:"data"`
	Timestamp int64  `json:"timestamp"`
	UserID    string `json:"user_id"`
}
type Subscription struct {
	ID      *string  `json:"id"`
	Topics  []string `json:"topics"`
	Command string   `json:"command"`
}

func (ctl *WebSocketController) Notification(c *gin.Context) {
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		service.SysLog.Panicln(err.Error())
		return
	}
	defer func() {
		if r := recover(); r != nil {
			service.SysLog.Errorln(r)
			if conn != nil {
				conn.Close()
			}
		}
	}()

	go ctl.Read(conn)
	go ctl.Write()

}
func (ctl *WebSocketController) Read(conn *websocket.Conn) {

	defer func() {
		if r := recover(); r != nil {
			service.SysLog.Errorln(r)
			if conn != nil {
				conn.Close()
				conn = nil
			}
		}
	}()
	for {

		if conn == nil {
			break
		}
		sub := &Subscription{}
		msgType, r, err := conn.NextReader()
		if err != nil {
			service.SysLog.Errorln(err.Error())
			conn.Close()
			conn = nil
			break
		}
		if msgType == -1 {
			conn.Close()
			conn = nil
			break
		}
		if msgType > -1 && err == nil {
			err := json.NewDecoder(r).Decode(sub)
			if err == nil {
				for _, topic := range sub.Topics {
					if ctl.Sockets[topic] == nil {
						ctl.Sockets[topic] = []*Connection{}
						ctl.PubSub.Register(topic, *sub.ID)
					}
					ctl.Sockets[topic] = append(ctl.Sockets[topic], &Connection{
						Socket: conn,
					})
				}
			}
		}
	}
}
func (ctl *WebSocketController) Write() {
	ctl.PubSub.Subscribe(func(msg *pubsub.Message) error {
		if msg != nil {
			for topic, connections := range ctl.Sockets {
				for _, connection := range connections {
					err := connection.Send(&Response{
						Topic:     topic,
						Data:      string(msg.Data),
						Timestamp: msg.Timestamp,
						UserID:    msg.UserID.String(),
					})
					if err != nil {
						connection.Socket.Close()
						delete(ctl.Sockets, topic)
					}
				}
			}
		}
		return nil
	})

}
