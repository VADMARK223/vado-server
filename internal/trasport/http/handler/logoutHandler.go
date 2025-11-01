package handler

import (
	"net/http"
	"vado_server/internal/config/code"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	c.SetCookie(
		code.JwtVado, // имя cookie
		"",           // пустое значение
		-1,           // MaxAge < 0 = удалить
		"/",          // путь
		"",           // домен
		false,        // secure (поставь true если HTTPS)
		true,         // httpOnly
	)

	c.Set(code.IsAuth, false)
	c.Redirect(http.StatusFound, "/")
}
