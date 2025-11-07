package handler

import (
	"fmt"
	"net/http"
	"time"
	"vado_server/internal/config/code"
	"vado_server/internal/domain/auth"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ShowIndex(secret string, log *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		td, _ := c.Get(code.TemplateData)
		data := td.(gin.H)

		tokenStr, errCookie := c.Cookie(code.JwtVado)
		log.Debugw("Index", "token", tokenStr, "err", errCookie)
		if errCookie != nil {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"Error": fmt.Sprintf("Error get token: %s", errCookie.Error()),
			})
			return
		}
		if errCookie == nil {
			claims, err := auth.ParseToken(tokenStr, secret)
			if err == nil {
				expTime := claims.ExpiresAt.Time
				remaining := time.Until(expTime).Truncate(time.Second)
				data[code.TokenExpireAt] = expTime.Format("02.01.2006 15:04:05")
				data[code.TokenRemaining] = remaining.String()
			} else {
				data[code.TokenExpireAt] = "токен невалиден"
			}
		} else {
			data[code.TokenExpireAt] = "токена нет"
		}

		c.HTML(http.StatusOK, "index.html", data)
	}

}
