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
		uidVal, exists := c.Get(code.UserId)
		if !exists {
			c.Next()
			return
		}

		if _, exists := c.Get(code.CurrentUser); exists {
			c.Next()
			return
		}

		userID := uidVal.(uint)

		if cached, ok := cache.Get(userID); ok {
			c.Set(code.CurrentUser, cached.(user.User))
			c.Next()
			return
		}

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
