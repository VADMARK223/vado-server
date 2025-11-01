package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowError(c *gin.Context, message string, err string) {
	c.HTML(http.StatusInternalServerError, "error.html", gin.H{
		"Message": message,
		"Error":   err,
	})
}
