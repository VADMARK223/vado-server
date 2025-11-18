package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	Conn   *websocket.Conn
	Hub    *Hub
	Send   chan []byte
	log    *zap.SugaredLogger
	UserID uint
}

type ClientPacket struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ServerMessage struct {
	Type     string `json:"type"` // "message"
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
	Text     string `json:"text"`
	Time     int64  `json:"time"`
}

func NewClient(conn *websocket.Conn, hub *Hub, userID uint, log *zap.SugaredLogger) *Client {
	client := &Client{
		Conn:   conn,
		Hub:    hub,
		Send:   make(chan []byte, 256),
		UserID: userID,
		log:    log,
	}
	return client
}

// IncomingLoop читает от клиента
func (c *Client) IncomingLoop() {
	defer func() {
		c.Hub.Unregister <- c
		_ = c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	_ = c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("WS read error:", err)
			break
		}

		var packet ClientPacket
		if errUnmarshal := json.Unmarshal(message, &packet); errUnmarshal != nil {
			log.Println("WS json error:", errUnmarshal)
			continue
		}

		c.log.Infow("Received message", "packet", packet)

		switch packet.Type {
		case "message":
			serverMsg := ServerMessage{
				Type:     "message",
				UserID:   c.UserID,
				Username: "test",
				Text:     packet.Text,
				Time:     time.Now().Unix(),
			}

			serverMsgBytes, _ := json.Marshal(serverMsg)

			c.Hub.Broadcast <- serverMsgBytes
		default:
			c.log.Warnw("Unknown message type", "type", packet.Type)
		}
	}
}

// OutgoingLoop вытаскивает из канала сообщения, которые присылает менеджер сообщений.
func (c *Client) OutgoingLoop() {
	ticker := time.NewTicker(30 * time.Second)

	defer func() {
		ticker.Stop()
		_ = c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Println("OutgoingLoop", "Channel closed:", err)
				}
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			_, err = w.Write(message)
			if err != nil {
				return
			}
			_ = w.Close()

		case <-ticker.C:
			// отправляем Ping
			_ = c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
