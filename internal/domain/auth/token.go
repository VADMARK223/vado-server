package auth

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/k0kubun/pp"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

func CreateTokenPair(userID uint, role string, accessTTL, refreshTTL time.Duration, secret string) (*TokenPair, error) {
	access, err := CreateToken(userID, role, accessTTL, secret)
	if err != nil {
		return nil, err
	}

	refresh, err := CreateToken(userID, role, refreshTTL, secret)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func CreateToken(userID uint, role string, ttl time.Duration, secret string) (string, error) {
	now := time.Now()
	log.Println("ROLE:", role)
	claims := CustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(int(userID)),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now), // Делает токен “активным не раньше определённого момента”
			Issuer:    "vado-server",
		},
	}

	_, _ = pp.Println("claims: ", claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenStr string, secret string) (*CustomClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("token is empty")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
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
