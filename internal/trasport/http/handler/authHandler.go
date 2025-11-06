package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"vado_server/internal/config/code"
	"vado_server/internal/config/route"
	"vado_server/internal/domain/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service      *user.Service
	secret       string
	tokenTTLSecs int
}

func NewAuthHandler(service *user.Service, secret string, tokenTTL string) *AuthHandler {
	tokenTTLSecs, _ := strconv.Atoi(tokenTTL)

	return &AuthHandler{
		service:      service,
		secret:       secret,
		tokenTTLSecs: tokenTTLSecs,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	_, token, _, err := h.service.Login(username, password, h.secret)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": err.Error()})
	}

	cookie := &http.Cookie{
		Name:     code.JwtVado,
		Value:    token,
		Path:     "/",
		HttpOnly: true,  // Нельзя прочитать из JS (document.cookie)
		Secure:   false, // Cookie отправляется даже по HTTP (Надо поменять в production)
		SameSite: http.SameSiteLaxMode,
		MaxAge:   h.tokenTTLSecs,
	}
	c.SetCookieData(cookie)
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
