package handlers

import (
	"fmt"
	"net/http"
	"vado_server/internal/appcontext"
	"vado_server/internal/models"

	"github.com/gin-gonic/gin"
)

func ShowRoles(appCtx *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var roles []models.Role
		if err := appCtx.DB.Find(&roles).Error; err != nil {
			appCtx.Log.Errorw("failed to get roles", "error", err)
			c.String(http.StatusInternalServerError, fmt.Sprintf("Ошибка получения ролей: %s", err.Error()))
			return
		}

		fmt.Println(len(roles))

		c.HTML(http.StatusOK, "roles.html", gin.H{
			"Roles": roles,
		})
	}
}
