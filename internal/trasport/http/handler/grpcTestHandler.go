package handler

import (
	"fmt"
	"net/http"
	"vado_server/internal/config/code"

	"github.com/gin-gonic/gin"
)

type GrpcTestHandler struct{}

func NewGrpcTestHandler() *GrpcTestHandler {
	return &GrpcTestHandler{}
}

func (h *GrpcTestHandler) ShowTestPage(c *gin.Context) {
	isAuth, ok := c.Get(code.IsAuth)
	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Message": fmt.Sprintf("Нет ключа (%s) в session", code.IsAuth),
			"Error":   fmt.Sprintf("Значение ключа: %v", isAuth),
		})
	}

	userID, ok := c.Get(code.UserId)

	if !ok {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Message": fmt.Sprintf("Нет ключа (%s) в session", code.UserId),
			"Error":   fmt.Sprintf("Значение ключа: %v", userID),
		})
	}

	//c.HTML(http.StatusOK, "grpc-test.html", gin.H{
	c.HTML(http.StatusOK, "hello.html", gin.H{
		code.IsAuth: isAuth,
		code.UserId: userID,
	})
}
