package auth

import (
	"net/http"
	"vado_server/internal/config/code"

	"github.com/gin-gonic/gin"
)

func SetTokenCookies(c *gin.Context, tokens *TokenPair, refreshTTL int) {
	SetCookie(c, code.VadoToken, tokens.AccessToken, refreshTTL)
	SetCookie(c, code.VadoRefreshToken, tokens.RefreshToken, refreshTTL)
}

func SetCookie(c *gin.Context, name, value string, maxAge int) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     "/",
		Domain:   "",
		Secure:   false, // Cookie отправляется даже по HTTP (Надо поменять в production) Защита от MITM
		HttpOnly: true,  // Нельзя прочитать из JS (document.cookie) Защита от XSS
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookieData(cookie)
}

func ClearTokenCookies(c *gin.Context) {
	clearCookie(c, code.VadoToken)
	clearCookie(c, code.VadoRefreshToken)
}

func clearCookie(c *gin.Context, name string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
