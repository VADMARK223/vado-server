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

	if errorMsg != "" {
		data["Error"] = errorMsg
	}

	c.HTML(http.StatusOK, "users.html", data)
}

func PostUser(service *user.Service) func(c *gin.Context) {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		email := c.PostForm("email")

		if username == "" {
			renderUsersPage(c, service, "Name is required")
			return
		}

		if password == "" {
			renderUsersPage(c, service, "Password is required")
			return
		}

		if email == "" {
			renderUsersPage(c, service, "Email is required")
			return
		}

		err := service.CreateUser(user.DTO{Username: username, Email: email, Password: password})

		if err != nil {
			ShowError(c, "Error adding user with role", err.Error())
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
