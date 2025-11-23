package handler

import (
	"net/http"
	"strconv"
	"vado_server/internal/config/code"
	"vado_server/internal/domain/user"

	"github.com/gin-gonic/gin"
)

func ShowUsers(service *user.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		renderUsersPage(c, service, "")
	}
}

func renderUsersPage(c *gin.Context, service *user.Service, errorMsg string) {
	users, err := service.GetAll()
	if err != nil {
		ShowError(c, "Failed to load users", err.Error())
		return
	}

	data := tplWithCapture(c, "Users list")
	data["Users"] = users

	if errorMsg != "" {
		data["Error"] = errorMsg
	}

	c.HTML(http.StatusOK, "users.html", data)
}

func DeleteUser(service *user.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		parseUint, parseUintErr := strconv.ParseUint(c.Param("id"), 10, 32)
		if parseUintErr != nil {
			return
		}

		currentId, errGetUSerId := c.Get(code.UserId)
		if !errGetUSerId {
			c.JSON(500, gin.H{
				"error": "Not fount user id",
			})
			return
		}

		if currentId == uint(parseUint) {
			c.JSON(400, gin.H{
				"error": "Вы пытаетесь удалить свой собственный аккаунт",
			})
			return
		}

		role, errRole := c.Get(code.Role)
		if !errRole {
			c.JSON(500, gin.H{
				"error": "Not fount user role",
			})
			return
		}

		if role != "admin" {
			c.JSON(400, gin.H{
				"error": "У вас нет прав на удаление аккаунтов",
			})
			return
		}

		err := service.DeleteUser(uint(parseUint))
		if err != nil {
			ShowError(c, "Failed to delete user", err.Error())
			return
		}

		c.JSON(200, gin.H{})
	}
}
