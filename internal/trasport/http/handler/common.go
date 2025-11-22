package handler

import (
	"vado_server/internal/config/code"

	"github.com/gin-gonic/gin"
)

func tplWithCapture(c *gin.Context, capture string) gin.H {
	td, exists := c.Get(code.TemplateData)
	if !exists {
		panic("TemplateData not found")
	}
	data := td.(gin.H)
	data[code.Capture] = capture

	return data
}
