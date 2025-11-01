package handler

import (
	"net/http"
	"strconv"
	"vado_server/internal/config/code"
	"vado_server/internal/config/route"
	"vado_server/internal/domain/user"

	"github.com/gin-gonic/gin"
)

func ShowUsers(service *user.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		var users, err = service.GetAllUsersWithRoles()
		if err != nil {
			ShowError(c, "Failed to load users", err.Error())
			return
		}

		td, _ := c.Get(code.TemplateData)
		data := td.(gin.H)
		data["Users"] = users
		c.HTML(http.StatusOK, "users.html", data)
	}
}

func PostUser(service *user.Service) func(c *gin.Context) {
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

func DeleteUser(service *user.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		parseUint, parseUintErr := strconv.ParseUint(c.Param("id"), 10, 32)
		if parseUintErr != nil {
			return
		}

		err := service.DeleteUser(uint(parseUint))
		if err != nil {
			ShowError(c, "Failed to delete user", err.Error())
			return
		}

		c.Redirect(http.StatusSeeOther, route.Users)
	}
}
