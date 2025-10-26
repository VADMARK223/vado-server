package handler

import (
	"fmt"
	"net/http"
	"vado_server/internal/app/context"
	"vado_server/internal/constants/code"
	"vado_server/internal/domain/task"

	"github.com/gin-gonic/gin"
)

func ShowTasksPage(service *task.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get(code.UserId)
		if !ok {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Message": "Нет ключа в session",
				"Error":   fmt.Sprintf("Значение ключа: %v", userID),
			})
		}

		if userID == nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Message": "user_id is nil",
			})
		}

		tasks, err := service.GetAllByUser(userID.(uint))
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Message": "Не удалось загрузить задачи",
			})
			return
		}

		isAuth, _ := c.Get(code.IsAuth)
		c.HTML(http.StatusOK, "tasks.html", gin.H{
			code.IsAuth: isAuth,
			code.UserId: userID,
			"Tasks":     tasks,
		})
	}
}

func AddTask(appCtx *context.AppContext) gin.HandlerFunc {
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

		t := task.Task{
			Name:        name,
			Description: desc,
			Completed:   completed == "on",
			UserID:      sessionUserID.(uint),
		}
		if err := appCtx.DB.Create(&t).Error; err != nil {
			appCtx.Log.Errorw("failed to create task", "error", err)
			ShowError(c, "Ошибка добавления задачи", err.Error())
			return
		}

		c.Redirect(http.StatusSeeOther, "/tasks")
	}
}

func DeleteTask(appCtx *context.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := appCtx.DB.Delete(&task.Task{}, id).Error; err != nil {
			appCtx.Log.Errorw("failed to delete task", "error", err)
			c.String(http.StatusInternalServerError, "Ошибка удаления задачи")
			return
		}

		c.Redirect(http.StatusSeeOther, "/tasks")
	}
}
