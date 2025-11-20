package handler

import (
	"net/http"
	"strconv"
	"vado_server/internal/config/code"
	"vado_server/internal/domain/user"

	"github.com/gin-gonic/gin"
	"github.com/k0kubun/pp"
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

	td, _ := c.Get(code.TemplateData)
	data := td.(gin.H)
	data["Users"] = users

	pp.Println("=============")
	pp.Println(len(users))

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

		err := service.DeleteUser(uint(parseUint))
		if err != nil {
			ShowError(c, "Failed to delete user", err.Error())
			return
		}

		c.JSON(200, gin.H{"status": "deleted"})
	}
}
