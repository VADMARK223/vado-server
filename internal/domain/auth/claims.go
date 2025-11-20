package auth

import (
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func (c *CustomClaims) UserID() uint {
	sub := c.Subject
	res, err := strconv.Atoi(sub)
	if err != nil {
		return 0
	}
	return uint(res)
}
