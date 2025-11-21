package middleware

import (
	"vado_server/internal/config/code"
	"vado_server/internal/domain/user"

	"github.com/gin-gonic/gin"
)

func TemplateContext(c *gin.Context) {
	result := gin.H{
		code.Username: "Guest",
		code.Mode:     gin.Mode(),
	}

	if contextUser, ok := c.Get(code.CurrentUser); ok {
		u := contextUser.(user.User)

		result[code.UserId] = u.ID
		result[code.Login] = u.Login
		result[code.Username] = u.Username
		result[code.Role] = u.Role
		result[code.Email] = u.Email
		result[code.IsAdmin] = u.IsAdmin()
	}

	c.Set(code.TemplateData, result)

	c.Next()
}
