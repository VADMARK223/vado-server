package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"vado_server/internal/app"
	"vado_server/internal/config/code"
	"vado_server/internal/domain/task"
	"vado_server/internal/infra/persistence/gorm"

	"github.com/gin-gonic/gin"
)

func Tasks(service *task.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := tplWithCapture(c, "User tasks")

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

		appCtx.Log.Debugw("Add task", "sessionUserID", sessionUserID)
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

func UpdateTask(appCtx *app.Context, service *task.Service) gin.HandlerFunc {
	type reqBody struct {
		Completed bool `json:"completed"`
	}

	return func(c *gin.Context) {
		appCtx.Log.Debugw("Update task")
		taskID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var body reqBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
			return
		}

		userID := c.MustGet(code.UserId).(uint)

		errUpdate := service.UpdateCompleted(uint(taskID), userID, body.Completed)
		if errUpdate != nil {
			appCtx.Log.Errorw("failed to update task", "error", errUpdate)
			c.JSON(http.StatusBadRequest, gin.H{"error": errUpdate.Error()})
			return
		}

		appCtx.Log.Debugw("Update task", "taskID", taskID, "userID", userID)
		c.Redirect(http.StatusSeeOther, "/tasks")
	}
}
