package handler

import (
	"fmt"
	"net/http"
	"time"
	"vado_server/internal/config/code"
	"vado_server/internal/infra/token"

	"github.com/gin-gonic/gin"
)

func ShowIndex(provider *token.JWTProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := tplWithCapture(c, fmt.Sprintf("Vado portal (%s)", gin.Mode()))

		updateTokenInfo(c, data, provider)
		updateRefreshTokenInfo(c, data, provider)

		c.HTML(http.StatusOK, "index.html", data)
	}
}

func updateTokenInfo(c *gin.Context, data gin.H, provider *token.JWTProvider) {
	data[code.TokenStatus] = "✅"
	data[code.TokenExpireAt] = "-"

	tokenStr, errTokenCookie := c.Cookie(code.VadoToken)
	if errTokenCookie != nil {
		data[code.TokenStatus] = "❌" + errTokenCookie.Error()
		return
	}

	claims, err := provider.ParseToken(tokenStr)
	if err != nil {
		data[code.TokenStatus] = "❌" + err.Error()
		return
	}
	expTime := claims.ExpiresAt.Time
	remaining := time.Until(expTime).Truncate(time.Second)
	data[code.TokenExpireAt] = fmt.Sprintf("%s (via %s)", expTime.Format("02.01.2006 15:04:05"), remaining.String())
}

func updateRefreshTokenInfo(c *gin.Context, data gin.H, provider *token.JWTProvider) {
	data[code.RefreshTokenStatus] = "✅"
	data[code.RefreshTokenExpireAt] = "-"

	tokenStr, errTokenCookie := c.Cookie(code.VadoRefreshToken)

	if errTokenCookie != nil {
		data[code.RefreshTokenStatus] = "❌" + errTokenCookie.Error()
		return
	}

	claims, err := provider.ParseToken(tokenStr)
	if err != nil {
		data[code.RefreshTokenStatus] = "❌" + err.Error()
		return
	}

	expTime := claims.ExpiresAt.Time
	remaining := time.Until(expTime).Truncate(time.Second)
	data[code.RefreshTokenExpireAt] = fmt.Sprintf("%s (via %s)", expTime.Format("02.01.2006 15:04:05"), remaining.String())
}
