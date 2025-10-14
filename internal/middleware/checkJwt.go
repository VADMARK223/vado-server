package middleware

import (
	"errors"
	"fmt"
	"time"
	"vado_server/internal/auth"
	"vado_server/internal/constants/code"
	"vado_server/internal/constants/env"
	"vado_server/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

/*func CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(code.JwtVado)
		if err != nil || tokenStr == "" {
			c.Set(code.IsAuth, false)
			c.Next()
			return
		}

		token, err := jwt.Parse(tokenStr,
			func(token *jwt.Token) (interface{}, error) {
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
			if exp, ok := claims[code.Exp].(float64); ok && int64(exp) < time.Now().Unix() {
				c.Set(code.IsAuth, false)
				c.Next()
				return
			}

			if userID, ok := claims[code.UserId]; ok {
				c.Set(code.UserId, userID)
			}
		}

		c.Set(code.IsAuth, true)
		c.Next()
	}
}*/

/*func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("token is null")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(util.GetEnv(env.JwtSecret)), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("token not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("not OK")
	}

	return claims, nil
}*/

func CheckJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(code.JwtVado)
		if err != nil || tokenStr == "" {
			c.Set(code.IsAuth, false)
			c.Next()
			return
		}

		claims, err := ParseToken(tokenStr)
		if err != nil {
			c.Set(code.IsAuth, false)
			c.Next()
			return
		}

		// Проверка срока действия токена
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			c.Set(code.IsAuth, false)
			c.Next()
			return
		}

		// Всё ок — записываем userID и флаг
		c.Set(code.IsAuth, true)
		c.Set(code.UserId, claims.UserID)
		c.Set("roles", claims.Roles)

		c.Next()
	}
}

func ParseToken(tokenStr string) (*auth.CustomClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("token is empty")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &auth.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(util.GetEnv(env.JwtSecret)), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	claims, ok := token.Claims.(*auth.CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
