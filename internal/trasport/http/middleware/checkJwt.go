package middleware

import (
	"time"
	"vado_server/internal/config/code"
	"vado_server/internal/domain/auth"

	"github.com/gin-gonic/gin"
)

func CheckJWT(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(code.JwtVado)
		if err != nil || tokenStr == "" {
			setNotAuth(c)
			c.Next()
			return
		}

		claims, err := auth.ParseToken(tokenStr, secret)
		if err != nil {
			setNotAuth(c)
			c.Next()
			return
		}

		// Проверка срока действия токена
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			setNotAuth(c)
			c.Next()
			return
		}

		// Всё ок — записываем userID и флаг
		c.Set(code.IsAuth, true)
		c.Set(code.UserId, claims.UserID)
		c.Set("roles", claims.Roles)

		c.Next()
	}
}

func setNotAuth(c *gin.Context) {
	c.Set(code.IsAuth, false)
	c.Set(code.UserId, "Guest")
}
