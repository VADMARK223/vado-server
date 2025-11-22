package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", tplWithCapture(c, "Sign in"))
}
