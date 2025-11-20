package handler

import (
	"fmt"
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

		updateTokenInfo(c, data, secret)
		updateRefreshTokenInfo(c, data, secret)

		c.HTML(http.StatusOK, "index.html", data)
	}
}

func updateTokenInfo(c *gin.Context, data gin.H, secret string) {
	data[code.TokenStatus] = "✅"
	data[code.TokenExpireAt] = "-"
	data[code.Role] = "-"

	tokenStr, errTokenCookie := c.Cookie(code.VadoToken)
	if errTokenCookie != nil {
		data[code.TokenStatus] = "❌" + errTokenCookie.Error()
		return
	}

	claims, err := auth.ParseToken(tokenStr, secret)
	if err != nil {
		data[code.TokenStatus] = "❌" + err.Error()
		return
	}
	expTime := claims.ExpiresAt.Time
	remaining := time.Until(expTime).Truncate(time.Second)
	data[code.TokenExpireAt] = fmt.Sprintf("%s (via %s)", expTime.Format("02.01.2006 15:04:05"), remaining.String())
	data[code.Role] = claims.Role
}

func updateRefreshTokenInfo(c *gin.Context, data gin.H, secret string) {
	data[code.RefreshTokenStatus] = "✅"
	data[code.RefreshTokenExpireAt] = "-"

	tokenStr, errTokenCookie := c.Cookie(code.VadoRefreshToken)

	if errTokenCookie != nil {
		data[code.RefreshTokenStatus] = "❌" + errTokenCookie.Error()
		return
	}

	claims, err := auth.ParseToken(tokenStr, secret)
	if err != nil {
		data[code.RefreshTokenStatus] = "❌" + err.Error()
		return
	}

	expTime := claims.ExpiresAt.Time
	remaining := time.Until(expTime).Truncate(time.Second)
	data[code.RefreshTokenExpireAt] = fmt.Sprintf("%s (via %s)", expTime.Format("02.01.2006 15:04:05"), remaining.String())
}
