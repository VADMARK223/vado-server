package middleware

import (
	"time"
	"vado_server/internal/config/code"
	"vado_server/internal/domain/auth"
	"vado_server/internal/infra/token"

	"github.com/gin-gonic/gin"
)

func CheckJWT(provider *token.JWTProvider, refreshTTL int) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(code.VadoToken)
		if err != nil || tokenStr == "" {
			tryRefresh(c, refreshTTL, provider)
			return
		}

		claims, err := provider.ParseToken(tokenStr)
		if err != nil {
			tryRefresh(c, refreshTTL, provider)
			return
		}

		// Проверка срока действия токена
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			tryRefresh(c, refreshTTL, provider)
			return
		}

		setAuth(c, claims.UserID(), claims.Role)
	}
}

func tryRefresh(c *gin.Context, refRefreshTokenTTL int, provider *token.JWTProvider) {
	refreshStr, err := c.Cookie(code.VadoRefreshToken)
	if err != nil || refreshStr == "" {
		c.Next()
		return
	}

	refreshClaims, err := provider.ParseToken(refreshStr)
	if err != nil || (refreshClaims.ExpiresAt != nil && refreshClaims.ExpiresAt.Time.Before(time.Now())) {
		c.Next()
		return
	}

	newAccess, err := provider.CreateToken(refreshClaims.UserID(), refreshClaims.Role, true)
	if err != nil {
		c.Next()
		return
	}

	auth.SetCookie(c, code.VadoToken, newAccess, refRefreshTokenTTL)

	setAuth(c, refreshClaims.UserID(), refreshClaims.Role)
}

func setAuth(c *gin.Context, id uint, role string) {
	c.Set(code.UserId, id)
	c.Set(code.Role, role)
	c.Next()
}
