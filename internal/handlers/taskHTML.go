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

func AddTask(appCtx *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		desc := c.PostForm("description")

		if name == "" {
			c.String(http.StatusBadRequest, "Название задачи обязательно")
			return
		}

		task := models.Task{
			Name:        name,
			Description: desc,
			UserID:      1, // TODO: for test
		}
		if err := appCtx.DB.Create(&task).Error; err != nil {
			appCtx.Log.Errorw("failed to create task", "error", err)
			c.String(http.StatusInternalServerError, "Ошибка добавления задачи")
			return
		}

		// После добавления возвращаем на список
		c.Redirect(http.StatusSeeOther, "/tasks")
	}
}
