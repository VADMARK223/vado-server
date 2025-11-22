package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Grpc(c *gin.Context) {
	c.HTML(http.StatusOK, "grpc-test.html", tplWithCapture(c, "Test gRPC"))
}
