package handler

import (
	"net/http"
	"time"
	"vado_server/internal/config/code"
	"vado_server/internal/domain/auth"

	"github.com/gin-gonic/gin"
)

func ShowIndex(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		td, _ := c.Get(code.TemplateData)
		data := td.(gin.H)

		tokenStr, errTokenCookie := c.Cookie(code.VadoToken)
		if errTokenCookie == nil {
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

		refreshTokenStr, errRefreshTokenCookie := c.Cookie(code.VadoRefreshToken)
		if errRefreshTokenCookie == nil {
			claims, err := auth.ParseToken(refreshTokenStr, secret)
			if err == nil {
				expTime := claims.ExpiresAt.Time
				remaining := time.Until(expTime).Truncate(time.Second)
				data[code.RefreshTokenExpireAt] = expTime.Format("02.01.2006 15:04:05")
				data[code.RefreshTokenRemaining] = remaining.String()
			} else {
				data[code.RefreshTokenExpireAt] = "refresh токен невалиден"
			}
		} else {
			data[code.RefreshTokenExpireAt] = "refresh токена нет"
		}

		c.HTML(http.StatusOK, "index.html", data)
	}

}
