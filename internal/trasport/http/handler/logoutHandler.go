package handler

import (
	"net/http"
	"vado_server/internal/config/code"
	"vado_server/internal/domain/auth"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	auth.ClearTokenCookies(c)
	c.Set(code.IsAuth, false)
	c.Redirect(http.StatusFound, "/")
}
