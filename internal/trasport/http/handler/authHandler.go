package handler

import (
	"fmt"
	"net/http"
	"vado_server/internal/config/code"
	"vado_server/internal/config/route"
	"vado_server/internal/domain/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *user.Service
}

func NewAuthHandler(service *user.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	_, token, _, err := h.service.Login(username, password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": err.Error()})
	}

	c.SetCookie(code.JwtVado,
		token,
		3600*24,
		"/",
		"",
		false,
		true)
	session := sessions.Default(c)

	redirectTo := session.Get(code.RedirectTo)
	if redirectTo == nil {
		redirectTo = route.Index
	} else {
		session.Delete(code.RedirectTo)
	}

	_ = session.Save()

	c.Redirect(http.StatusFound, redirectTo.(string))
}

func ShowLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	}
}

func ShowRegisterPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	}
}

func PerformRegister(service *user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")

		err := service.CreateUser(user.DTO{Username: username, Email: email, Password: password})

		if err != nil {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"Error": fmt.Sprintf("Ошибка при регистрации: %s", err.Error()),
			})
			return
		}

		c.Redirect(http.StatusFound, "/login")
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie(
			code.JwtVado, // имя cookie
			"",           // пустое значение
			-1,           // MaxAge < 0 = удалить
			"/",          // путь
			"",           // домен
			false,        // secure (поставь true если HTTPS)
			true,         // httpOnly
		)

		c.Set(code.IsAuth, false)

		c.Redirect(http.StatusFound, "/")
	}
}
