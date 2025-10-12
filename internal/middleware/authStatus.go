package middleware

import (
	"vado_server/internal/constants/code"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthStatusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get(code.UserId)
		isAuth := userID != nil
		c.Set(code.IsAuth, isAuth)
		c.Next()
	}
}
