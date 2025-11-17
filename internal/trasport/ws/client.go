package ws

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	Conn *websocket.Conn
	Hub  *Hub
	Send chan []byte
	log  *zap.SugaredLogger
}

func NewClient(conn *websocket.Conn, hub *Hub, log *zap.SugaredLogger) *Client {
	client := &Client{
		Conn: conn,
		Hub:  hub,
		Send: make(chan []byte, 256),
		log:  log,
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

		c.log.Infow("Received message", "message", string(message))
		c.Hub.Broadcast <- message
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

			username := "VADMARK"
			_, err = w.Write([]byte(fmt.Sprintf("%s: %s", username, message)))
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
