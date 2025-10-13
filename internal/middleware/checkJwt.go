package middleware

import (
	"time"
	"vado_server/internal/constants/code"
	"vado_server/internal/constants/env"
	"vado_server/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(code.JwtVado)
		if err != nil || tokenStr == "" {
			c.Set(code.IsAuth, false)
			c.Next()
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(util.GetEnv(env.JwtSecret)), nil
		})

		if err != nil || !token.Valid {
			c.Set(code.IsAuth, false)
			c.Next()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
				c.Set(code.IsAuth, false)
				c.Next()
				return
			}

			if userID, ok := claims["user_id"]; ok {
				c.Set("user_id", userID)
			}
		}

		c.Set(code.IsAuth, true)
		c.Next()
	}
}
