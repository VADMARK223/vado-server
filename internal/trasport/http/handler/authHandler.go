package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"vado_server/internal/config/code"
	"vado_server/internal/config/route"
	"vado_server/internal/domain/auth"
	"vado_server/internal/domain/user"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	service             *user.Service
	secret              string
	tokenTTLSecs        int
	refreshTokenTTLSecs int
	log                 *zap.SugaredLogger
}

func NewAuthHandler(service *user.Service, secret string, tokenTTL string, refreshTokenTTL string, log *zap.SugaredLogger) *AuthHandler {
	tokenTTLSecs, _ := strconv.Atoi(tokenTTL)
	refreshTokenTTLSecs, _ := strconv.Atoi(refreshTokenTTL)

	return &AuthHandler{
		service:             service,
		secret:              secret,
		tokenTTLSecs:        tokenTTLSecs,
		refreshTokenTTLSecs: refreshTokenTTLSecs,
		log:                 log,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	_, tokens, err := h.service.Login(username, password)
	h.log.Debugw("Token from service", "tokens", tokens)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": err.Error()})
	}

	auth.SetTokenCookies(c, tokens, h.refreshTokenTTLSecs)

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
		role := c.PostForm(code.Role)
		color := c.PostForm(code.Color)

		err := service.CreateUser(user.DTO{Username: username, Email: email, Password: password, Role: user.Role(role), Color: color})

		if err != nil {
			c.HTML(http.StatusBadRequest, "register.html", gin.H{
				"Error": fmt.Sprintf("Ошибка при регистрации: %s", err.Error()),
			})
			return
		}

		c.Redirect(http.StatusFound, "/login")
	}
}
