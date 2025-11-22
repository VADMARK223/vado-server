package token

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"vado_server/internal/domain/auth"

	"github.com/golang-jwt/jwt/v5"
)

type JWTProvider struct {
	secret     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTProvider(secret string, accessTTL, refreshTTL time.Duration) *JWTProvider {
	return &JWTProvider{secret: secret, accessTTL: accessTTL, refreshTTL: refreshTTL}
}

func (j *JWTProvider) CreateTokenPair(userID uint, role string) (*auth.TokenPair, error) {
	access, err := j.CreateToken(userID, role, true)
	if err != nil {
		return nil, err
	}

	refresh, err := j.CreateToken(userID, role, false)
	if err != nil {
		return nil, err
	}

	return &auth.TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (j *JWTProvider) CreateToken(userID uint, role string, accessToken bool) (string, error) {
	now := time.Now()

	var duration time.Duration
	if accessToken {
		duration = j.accessTTL
	} else {
		duration = j.refreshTTL
	}

	claims := auth.CustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(int(userID)),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now), // Делает токен “активным не раньше определённого момента”
			Issuer:    "vado-server",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWTProvider) ParseToken(tokenStr string) (*auth.CustomClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("token is empty")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &auth.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(j.secret), nil
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
