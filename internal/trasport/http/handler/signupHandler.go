package handler

import (
	"net/http"
	"vado_server/internal/config/code"

	"github.com/gin-gonic/gin"
)

func ShowSignup(c *gin.Context) {
	data := tplWithCapture(c, "Sign up")
	if errVal, exists := c.Get(code.Error); exists {
		data[code.Error] = errVal
	}
	c.HTML(http.StatusOK, "register.html", data)
}
