package handler

import (
	"net/http"
	"strings"
	"vado_server/internal/config/code"
	"vado_server/internal/config/route"
	"vado_server/internal/domain/auth"
	"vado_server/internal/domain/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service    *user.Service
	secret     string
	refreshTTL int
	log        *zap.SugaredLogger
}

func NewAuthHandler(service *user.Service, secret string, refreshTTL int, log *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{
		service:    service,
		secret:     secret,
		refreshTTL: refreshTTL,
		log:        log,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	login := c.PostForm("login")
	password := c.PostForm("password")

	_, tokens, err := h.service.Login(login, password)
	if err != nil {
		c.Set(code.Error, strings.ToUpper(err.Error()[:1])+err.Error()[1:])
		ShowLogin(c)
		return
	}

	auth.SetTokenCookies(c, tokens, h.refreshTTL)

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
		login := c.PostForm(code.Login)
		email := c.PostForm(code.Email)
		password := c.PostForm(code.Password)
		role := c.PostForm(code.Role)
		color := c.PostForm(code.Color)
		username := c.PostForm(code.Username)

		err := service.CreateUser(user.DTO{Login: login, Email: email, Password: password, Role: user.Role(role), Color: color, Username: username})

		if err != nil {
			c.Set(code.Error, strings.ToUpper(err.Error()[:1])+err.Error()[1:])
			ShowSignup(c)
			return
		}

		c.Redirect(http.StatusFound, "/login")
	}
}
