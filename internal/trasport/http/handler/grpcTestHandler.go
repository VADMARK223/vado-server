package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GrpcTestHandler struct{}

func NewGrpcTestHandler() *GrpcTestHandler {
	return &GrpcTestHandler{}
}

func (h *GrpcTestHandler) ShowTestPage(c *gin.Context) {
	c.HTML(http.StatusOK, "grpc-test.html", gin.H{
		"title": "gRPC test",
	})
}
