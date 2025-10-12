package handlers

import (
	"net/http"
	"vado_server/internal/constants/code"

	"github.com/gin-gonic/gin"
)

func ShowIndex(c *gin.Context) {
	isAuth, _ := c.Get(code.IsAuth)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Mode":      gin.Mode(),
		code.IsAuth: isAuth,
	})
}
