package handler

import (
	"log"
	"net/http"
	"vado_server/internal/config/code"
	"vado_server/internal/trasport/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func ShowChat() func(c *gin.Context) {
	return func(c *gin.Context) {
		td, _ := c.Get(code.TemplateData)
		data := td.(gin.H)
		c.HTML(http.StatusOK, "chat.html", data)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Разрешаем любые origins (можно ужесточить)
	CheckOrigin: func(r *http.Request) bool { return true },
}

func ChatHandler(hub *ws.Hub, logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}

		client := ws.NewClient(conn, hub, logger)
		hub.Register <- client

		go client.OutgoingLoop()
		go client.IncomingLoop()
	}
}
