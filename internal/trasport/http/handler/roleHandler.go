package handler

import (
	"net/http"
	"vado_server/internal/config/code"
	"vado_server/internal/domain/role"

	"github.com/gin-gonic/gin"
)

func Roles(service *role.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, err := service.GetAll()
		if err != nil {
			ShowError(c, "Failed to load roles", err.Error())
			return
		}

		td, _ := c.Get(code.TemplateData)
		data := td.(gin.H)
		data["Roles"] = roles
		c.HTML(http.StatusOK, "roles.html", data)
	}
}
