package handler

import (
	"net/http"
	"vado_server/internal/config/code"
	"vado_server/internal/domain/auth"
	"vado_server/internal/domain/user"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func MeHandler(log *zap.SugaredLogger, secret string, svc *user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, errTokenCookie := c.Cookie(code.VadoToken)
		if errTokenCookie != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": errTokenCookie.Error()})
			return
		}

		claims, err := auth.ParseToken(tokenStr, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"Error": errTokenCookie.Error()})
			return
		}

		userId := claims.UserID()
		u, errGetUser := svc.GetByID(userId)
		if errGetUser != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"Error": errGetUser.Error()})
			return
		}
		log.Infow("MeHandler:", "email", u.Email)
		c.JSON(200, gin.H{"email": u.Email})
	}
}
