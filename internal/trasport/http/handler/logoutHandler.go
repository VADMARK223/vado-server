package handler

import (
	"net/http"
	"vado_server/internal/domain/auth"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	auth.ClearTokenCookies(c)
	c.Redirect(http.StatusFound, "/")
}
