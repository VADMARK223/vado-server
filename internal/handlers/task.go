package handlers

import (
	"net/http"
	"vado_server/internal/appcontext"
	"vado_server/internal/models"

	"github.com/gin-gonic/gin"
)

func GetTasks(ctx *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tasks []models.Task
		if err := ctx.DB.Find(&tasks).Error; err != nil {
			ctx.Log.Errorw("Failed to get tasks", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tasks"})
			return
		}
		c.JSON(http.StatusOK, tasks)
	}
}
