package handler

import (
	"fmt"
	"net/http"
	"time"
	"vado_server/internal/app/context"
	"vado_server/internal/auth"
	"vado_server/internal/config/code"
	"vado_server/internal/config/route"
	"vado_server/internal/domain/user"
	user2 "vado_server/internal/infra/persistence/user"
	auth3 "vado_server/internal/trasport/grpc/auth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ShowLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	}
}

func PerformLogin(appCtx *context.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		//var u user.User
		var u user2.Entity
		if err := appCtx.DB.Where("username = ?", username).First(&u).Error; err != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "Пользователь не найден"})
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "Неверный пароль"})
			return
		}

		token, err := auth.CreateToken(u.ID, []string{"user"}, time.Minute*auth3.TokenAliveMinutes)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "error.html", gin.H{
				"Message": "Ошибка генерации токена",
				"Error":   err.Error(),
			})
			return
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
