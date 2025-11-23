package handler

import (
	"net/http"
	"vado_server/internal/config/code"

	"github.com/gin-gonic/gin"
)

func Grpc(c *gin.Context) {
	token, _ := c.Cookie(code.VadoToken)
	data := tplWithCapture(c, "Test gRPC")
	data[code.VadoToken] = token
	c.HTML(http.StatusOK, "grpc-test.html", data)
}
