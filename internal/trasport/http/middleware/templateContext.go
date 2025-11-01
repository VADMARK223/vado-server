package middleware

import (
	"fmt"
	"net/http"
	"vado_server/internal/config/code"

	"github.com/gin-gonic/gin"
)

func TemplateContext(c *gin.Context) {
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

	c.Set(code.TemplateData, gin.H{
		code.IsAuth: isAuth,
		code.UserId: userID,
		code.Mode:   gin.Mode(),
	})

	c.Next()
}
