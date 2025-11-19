package middleware

import (
	"strconv"
	"time"
	"vado_server/internal/config/code"
	"vado_server/internal/domain/auth"

	"github.com/gin-gonic/gin"
)

func CheckJWT(secret, tokenTTL, refreshTokenTTL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenTTLSecs, _ := strconv.Atoi(tokenTTL)
		refreshTokenTTLSecs, _ := strconv.Atoi(refreshTokenTTL)
		tokenStr, err := c.Cookie(code.VadoToken)
		if err != nil || tokenStr == "" {
			tryRefresh(c, secret, tokenTTLSecs, refreshTokenTTLSecs)
			return
		}

		claims, err := auth.ParseToken(tokenStr, secret)
		if err != nil {
			tryRefresh(c, secret, tokenTTLSecs, refreshTokenTTLSecs)
			return
		}

		// Проверка срока действия токена
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			tryRefresh(c, secret, tokenTTLSecs, refreshTokenTTLSecs)
			return
		}

		setAuth(c, claims.UserID(), claims.Roles)
	}
}

func tryRefresh(c *gin.Context, secret string, tokenTTL, refRefreshTokenTTL int) {
	refreshStr, err := c.Cookie(code.VadoRefreshToken)
	if err != nil || refreshStr == "" {
		setNotAuth(c)
		return
	}

	refreshClaims, err := auth.ParseToken(refreshStr, secret)
	if err != nil || (refreshClaims.ExpiresAt != nil && refreshClaims.ExpiresAt.Time.Before(time.Now())) {
		setNotAuth(c)
		return
	}

	newAccess, err := auth.CreateToken(refreshClaims.UserID(), refreshClaims.Roles, time.Second*time.Duration(tokenTTL), secret)
	if err != nil {
		setNotAuth(c)
		return
	}

	auth.SetCookie(c, code.VadoToken, newAccess, refRefreshTokenTTL)

	setAuth(c, refreshClaims.UserID(), refreshClaims.Roles)
}

func setAuth(c *gin.Context, id uint, roles []uint) {
	c.Set(code.IsAuth, true)
	c.Set(code.UserId, id)
	c.Set("roles", roles)
	c.Next()
}

func setNotAuth(c *gin.Context) {
	c.Set(code.IsAuth, false)
	c.Set(code.UserId, "Guest")
	c.Next()
}
