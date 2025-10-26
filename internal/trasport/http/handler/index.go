package handler

import (
	"net/http"
	"vado_server/internal/constants/code"

	"github.com/gin-gonic/gin"
)

func ShowIndex(c *gin.Context) {
	isAuth, _ := c.Get(code.IsAuth)
	userId, _ := c.Get(code.UserId)
	c.HTML(http.StatusOK, "index.html", gin.H{
		code.IsAuth: isAuth,
		code.UserId: userId,
		code.Mode:   gin.Mode(),
	})
}
