package handlers

import (
	"net/http"
	"vado_server/internal/appcontext"
	"vado_server/internal/constants/code"
	"vado_server/internal/models"
	"vado_server/internal/services"

	"github.com/gin-gonic/gin"
)

func ShowUsers(appCtx *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		if err := appCtx.DB.Preload("Roles").Find(&users).Error; err != nil {
			ShowError(c, "Не удалось загрузить пользователей", err.Error())
			return
		}

		isAuth, _ := c.Get(code.IsAuth)
		userId, _ := c.Get(code.UserId)
		c.HTML(http.StatusOK, "users.html", gin.H{
			code.IsAuth: isAuth,
			code.UserId: userId,
			"Users":     users,
		})
	}
}

func AddUser(service *services.UserService) func(c *gin.Context) {
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

		userDto := models.UserDTO{
			Username: username,
			Password: password,
			Email:    email,
		}

		err := service.CreateUser(userDto)

		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Message": "Ошибка добавления пользователя c ролью",
				"Error":   err.Error(),
			})
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
