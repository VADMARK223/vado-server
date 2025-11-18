package handler

import (
	"net/http"
	"vado_server/internal/config/code"
	"vado_server/internal/domain/auth"

	//"vado_server/internal/domain/auth"
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

func ServeSW(hub *ws.Hub, log *zap.SugaredLogger, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Query("token")
		log.Infow("ServeSW", "tokenStr", tokenStr)
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
			return
		}

		claims, err := auth.ParseToken(tokenStr, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Errorw("Upgrader error", "error", err)
			return
		}

		client := ws.NewClient(conn, hub, claims.UserID, log)
		hub.Register <- client

		go client.OutgoingLoop()
		go client.IncomingLoop()
	}
}
