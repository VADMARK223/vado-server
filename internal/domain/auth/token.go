package auth

import (
	"errors"
	"fmt"
	"time"
	"vado_server/internal/config/env"
	"vado_server/internal/util"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID uint     `json:"user_id"`
	Roles  []string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

func CreateToken(userID uint, roles []string, ttl time.Duration) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "vado-ping",
			Subject:   "access",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(util.GetEnv(env.JwtSecret)))
}

func ParseToken(tokenStr string) (*CustomClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("token is empty")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(util.GetEnv(env.JwtSecret)), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
