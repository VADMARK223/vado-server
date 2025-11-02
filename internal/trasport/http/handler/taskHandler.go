package handler

import (
	"fmt"
	"net/http"
	"vado_server/internal/app"
	"vado_server/internal/config/code"
	"vado_server/internal/domain/task"
	"vado_server/internal/infra/persistence/gorm"

	"github.com/gin-gonic/gin"
)

func Tasks(service *task.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		td, _ := c.Get(code.TemplateData)
		data := td.(gin.H)

		tasks, err := service.GetAllByUser(data[code.UserId].(uint))
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Message": "Не удалось загрузить задачи",
			})
			return
		}

		data["Tasks"] = tasks
		c.HTML(http.StatusOK, "tasks.html", data)
	}
}

func AddTask(appCtx *app.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.PostForm("name")
		desc := c.PostForm("description")
		completed := c.PostForm("completed")
		appCtx.Log.Debugw("Add task", "name", name, "desc", desc, "completed", completed)

		if name == "" {
			c.String(http.StatusBadRequest, "Название задачи обязательно")
			return
		}

		sessionUserID, ok := c.Get(code.UserId)
		if !ok {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Message": "Нет ключа в session",
				"Error":   fmt.Sprintf("Значение ключа: %v", code.UserId),
			})
		}

		t := gorm.TaskEntity{
			Name:        name,
			Description: desc,
			Completed:   completed == "on",
			UserID:      sessionUserID.(uint),
		}

		if err := appCtx.DB.Create(&t).Error; err != nil {
			appCtx.Log.Errorw("failed to create task", "error", err)
			ShowError(c, "Error adding task", err.Error())
			return
		}

		c.Redirect(http.StatusSeeOther, "/tasks")
	}
}

func DeleteTask(appCtx *app.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := appCtx.DB.Delete(&task.Task{}, id).Error; err != nil {
			appCtx.Log.Errorw("failed to delete task", "error", err)
			c.String(http.StatusInternalServerError, "Error deleting task")
			return
		}

		c.Redirect(http.StatusSeeOther, "/tasks")
	}
}
