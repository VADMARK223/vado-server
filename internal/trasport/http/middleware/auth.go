package middleware

import (
	"net/http"
	"time"
	"vado_server/internal/config/code"
	"vado_server/internal/domain/auth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAuth, ok := c.Get(code.IsAuth)
		if !ok || isAuth == false {
			session := sessions.Default(c)
			session.Set(code.RedirectTo, c.Request.URL.Path)
			_ = session.Save()
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Next()
	}
}

func CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(code.JwtVado)
		if err != nil || tokenStr == "" {
			c.Set(code.IsAuth, false)
			c.Next()
			return
		}

		claims, err := auth.ParseToken(tokenStr)
		if err != nil {
			c.Set(code.IsAuth, false)
			c.Next()
			return
		}

		// Проверка срока действия токена
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			c.Set(code.IsAuth, false)
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
