package handlers

import (
	"net/http"
	"vado_server/internal/appcontext"
	"vado_server/internal/models"

	"github.com/gin-gonic/gin"
)

func ShowUsers(appCtx *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		if err := appCtx.DB.Preload("Roles").Find(&users).Error; err != nil {
			ShowError(c, "Не удалось загрузить пользователей", err.Error())
			return
		}

		c.HTML(http.StatusOK, "users.html", gin.H{
			"Users": users,
		})
	}
}

func AddUser(appCtx *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		email := c.PostForm("email")

		if username == "" {
			c.String(http.StatusBadRequest, "Имя обязательно")
			return
		}

		if password == "" {
			c.String(http.StatusBadRequest, "Пароль обязателен")
			return
		}

		user := models.User{
			Username: username,
			Password: password,
			Email:    email,
		}
		if err := appCtx.DB.Create(&user).Error; err != nil {
			appCtx.Log.Errorw("failed to create user", "error", err)
			c.String(http.StatusInternalServerError, "Ошибка добавления пользователя")
			return
		}

		c.Redirect(http.StatusSeeOther, "/users")
	}
}

func DeleteUser(appCtx *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := appCtx.DB.Delete(&models.User{}, id).Error; err != nil {
			appCtx.Log.Errorw("failed to delete user", "error", err)
			c.String(http.StatusInternalServerError, "Ошибка удаления пользователя")
			return
		}

		c.Redirect(http.StatusSeeOther, "/users")
	}
}
