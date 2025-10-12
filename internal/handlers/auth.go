package handlers

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear() // удаляем все данные из сессии
		session.Options(sessions.Options{
			Path:     "/", // обязательно, иначе cookie не перезапишется
			MaxAge:   -1,  // удалить cookie сразу
			HttpOnly: true,
		})
		_ = session.Save()

		c.Redirect(http.StatusFound, "/")
	}
}
