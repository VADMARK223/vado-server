package handler

import (
	"net/http"
	"vado_server/internal/config/code"

	"github.com/gin-gonic/gin"
)

func ShowSignup(c *gin.Context) {
	td, _ := c.Get(code.TemplateData)
	data := td.(gin.H)
	c.HTML(http.StatusOK, "register.html", data)
}
