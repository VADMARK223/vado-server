package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowBook(c *gin.Context) {
	c.HTML(http.StatusOK, "book.html", tplWithCapture(c, "Golang book"))
}
