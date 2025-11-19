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

func ShowChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		td, _ := c.Get(code.TemplateData)
		data := td.(gin.H)

		tokenStr, errTokenCookie := c.Cookie(code.VadoToken)
		if errTokenCookie == nil && tokenStr != "" {
			data["ws_token"] = tokenStr
		} else {
			data["ws_token"] = ""
		}

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
		// Если токена нет в query параметре, пробуем взять из кук
		if tokenStr == "" {
			var err error
			tokenStr, err = c.Cookie(code.VadoToken)
			if err != nil || tokenStr == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
				return
			}
		}
		log.Infow("ServeSW", "tokenStr", tokenStr)

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

		client := ws.NewClient(conn, hub, claims.UserID(), log)
		hub.Register <- client

		go client.OutgoingLoop()
		go client.IncomingLoop()
	}
}
