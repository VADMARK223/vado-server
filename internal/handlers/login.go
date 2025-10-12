package handlers

import (
	"net/http"
	"vado_server/internal/appcontext"
	"vado_server/internal/constants/code"
	"vado_server/internal/constants/route"
	"vado_server/internal/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
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

		if user.Password != password {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "Неверный пароль"})
			return
		}

		/*if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "Неверный пароль"})
			return
		}*/

		session := sessions.Default(c)
		session.Set(code.UserId, user.ID)

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
