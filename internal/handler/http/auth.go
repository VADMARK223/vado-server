package http

import (
	"fmt"
	"net/http"
	"time"
	"vado_server/internal/appcontext"
	"vado_server/internal/auth"
	"vado_server/internal/constants/code"
	"vado_server/internal/constants/route"
	"vado_server/internal/models"
	"vado_server/internal/services"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func ShowLoginPage() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	}
}

func PerformLogin(appCtx *appcontext.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		var user models.User
		if err := appCtx.DB.Where("username = ?", username).First(&user).Error; err != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "Пользователь не найден"})
			return
		}

		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "Неверный пароль"})
			return
		}

		token, err := auth.CreateToken(user.ID, []string{"user"}, time.Minute*15)
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

func PerformRegister(service *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		email := c.PostForm("email")
		password := c.PostForm("password")

		err := service.CreateUser(models.UserDTO{Username: username, Email: email, Password: password})

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
