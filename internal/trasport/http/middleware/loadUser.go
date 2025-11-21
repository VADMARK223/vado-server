package middleware

import (
	"net/http"
	"time"
	"vado_server/internal/app"
	"vado_server/internal/config/code"
	"vado_server/internal/domain/user"

	"github.com/gin-gonic/gin"
)

func LoadUserContext(svc *user.Service, cache *app.LocalCache) gin.HandlerFunc {
	return func(c *gin.Context) {
		app.Dump("Load user context", nil)

		uidVal, exists := c.Get(code.UserId)
		if !exists {
			app.Dump("Пользователь не авторизован", nil)
			c.Next()
			return
		}

		if _, exists := c.Get(code.CurrentUser); exists {
			app.Dump("Пользователь уже загружен", nil)
			c.Next()
			return
		}

		userID := uidVal.(uint)

		if cached, ok := cache.Get(userID); ok {
			app.Dump("Достаем пользователя из кеша", nil)
			c.Set(code.CurrentUser, cached.(user.User))
			c.Next()
			return
		}

		app.Dump("Запрашиваем пользователя из БД", nil)
		u, err := svc.GetByID(userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
			return
		}

		cache.Set(userID, u, time.Minute*5)
		c.Set(code.CurrentUser, u)

		c.Next()
	}
}
