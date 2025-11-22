package handler

import (
	"net/http"
	"vado_server/internal/config/code"

	"github.com/gin-gonic/gin"
)

func ShowLogin(c *gin.Context) {
	data := tplWithCapture(c, "Sign in")
	if errVal, exists := c.Get(code.Error); exists {
		data[code.Error] = errVal
	}

	c.HTML(http.StatusOK, "login.html", data)
}
