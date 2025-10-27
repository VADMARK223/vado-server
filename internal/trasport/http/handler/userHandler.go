package handler

import (
	"net/http"
	"vado_server/internal/app/context"
	"vado_server/internal/config/code"
	"vado_server/internal/config/route"
	"vado_server/internal/domain/user"
	user2 "vado_server/internal/infra/persistence/gorm"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *user.Service
}

func (h *UserHandler) ShowUsers(c *gin.Context) {
	var users, err = h.service.GetAllUsersWithRoles()
	if err != nil {
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

func NewUserHandler(service *user.Service) *UserHandler {
	return &UserHandler{service: service}
}

func AddUser(service *user.Service) func(c *gin.Context) {
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

		userDto := user.DTO{
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
		c.Redirect(http.StatusSeeOther, route.Users)
	}
}

func DeleteUser(appCtx *context.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := appCtx.DB.Delete(&user2.UserEntity{}, id).Error; err != nil {
			appCtx.Log.Errorw("failed to delete user", "error", err)
			c.String(http.StatusInternalServerError, "Ошибка удаления пользователя")
			return
		}

		c.Redirect(http.StatusSeeOther, route.Users)
	}
}
