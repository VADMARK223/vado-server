package handlers

import (
	"net/http"
	"vado_server/internal/appcontext"
	"vado_server/internal/models"

	"github.com/gin-gonic/gin"
)

func ShowAdmin(appCtx *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		if err := appCtx.DB.Find(&users).Error; err != nil {
			appCtx.Log.Errorw("failed to get users", "error", err)
			c.String(http.StatusInternalServerError, "Ошибка получения пользователей")
			return
		}

		c.HTML(http.StatusOK, "admin.html", gin.H{
			"Users": users,
		})
	}
}

func AddUser(appCtx *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

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
		}
		if err := appCtx.DB.Create(&user).Error; err != nil {
			appCtx.Log.Errorw("failed to create user", "error", err)
			c.String(http.StatusInternalServerError, "Ошибка добавления пользователя")
			return
		}

		c.Redirect(http.StatusSeeOther, "/admin")
	}
}
