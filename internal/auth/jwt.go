package auth

import (
	"time"
	"vado_server/internal/constants/env"
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
			Issuer:    "vado-server",
			Subject:   "access",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(util.GetEnv(env.JwtSecret)))
}
