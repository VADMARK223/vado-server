package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ShowSignup(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", tplWithCapture(c, "Sign up"))
}
