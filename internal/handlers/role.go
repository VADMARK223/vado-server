package handlers

import (
	"net/http"
	"vado_server/internal/services"

	"github.com/gin-gonic/gin"
)

func ShowRoles(service *services.RoleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, err := service.GetAll()
		if err != nil {
			ShowError(c, "Не удалось загрузить роли", err.Error())
			return
		}

		c.HTML(http.StatusOK, "roles.html", gin.H{
			"Roles": roles,
		})
	}
}
