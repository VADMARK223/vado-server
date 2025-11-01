package middleware

import (
	"net/http"
	"vado_server/internal/config/code"
	"vado_server/internal/config/route"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CheckAuthAndRedirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAuth, ok := c.Get(code.IsAuth)
		if !ok || isAuth == false {
			session := sessions.Default(c)
			session.Set(code.RedirectTo, c.Request.URL.Path)
			_ = session.Save()
			c.Redirect(http.StatusFound, route.Login)
			c.Abort()
			return
		}

		c.Next()
	}
}
