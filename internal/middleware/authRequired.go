package middleware

import (
	"net/http"
	"vado_server/internal/constants/code"
	"vado_server/internal/constants/route"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AuthRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get(code.UserId)

		if userID == nil {
			session.Set(code.RedirectTo, c.Request.URL.Path)
			_ = session.Save()

			c.Redirect(http.StatusFound, route.Login)
			c.Abort()
			return
		}
		c.Next()
	}
}
