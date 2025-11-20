package auth

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type TokenProvider interface {
	CreateTokenPair(userID uint, role string) (*TokenPair, error)
	CreateToken(userID uint, role string) (string, error)
	ParseToken(token string) (*CustomClaims, error)
}
