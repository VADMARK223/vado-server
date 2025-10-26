package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MessageDTO struct {
	Text string `json:"text"`
}

func GetHello() gin.HandlerFunc {
	return func(c *gin.Context) {
		//message := MessageDTO{Text: "Hello World"}
		c.JSON(http.StatusInternalServerError, gin.H{
			"Text": "hello from server!",
		})
	}
}
