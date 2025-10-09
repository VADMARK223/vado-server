package handlers

import (
	"net/http"
	"vado_server/internal/appcontext"
	"vado_server/internal/models"

	"github.com/gin-gonic/gin"
)

func ShowTasks(appCtx *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tasks []models.Task
		if err := appCtx.DB.Find(&tasks).Error; err != nil {
			appCtx.Log.Errorw("failed to get tasks", "error", err)
			c.String(http.StatusInternalServerError, "Ошибка получения задач")
			return
		}

		c.HTML(http.StatusOK, "tasks.html", gin.H{
			"Tasks": tasks,
		})
	}
}
