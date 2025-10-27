package user

import (
	"errors"
	"time"
	"vado_server/internal/domain/auth"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo                Repository
	accessTokenDuration time.Duration
}

func NewService(repo Repository, accessTokenDuration time.Duration) *Service {
	return &Service{
		repo:                repo,
		accessTokenDuration: accessTokenDuration,
	}
}

func (s *Service) CreateUser(dto DTO) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	user := User{
		Username: dto.Username,
		Password: string(hash),
		Email:    dto.Email,
	}
	return s.repo.CreateUser(user)
}

func (s *Service) Login(username string, password string) (string, error) {
	u, errGetUser := s.repo.GetByUsername(username)
	if errGetUser != nil {
		return "", errors.New("пользователь не найден")
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		return "", errors.New("неверный пароль")
	}

	token, errToken := auth.CreateToken(u.ID, []string{"user"}, s.accessTokenDuration)
	if errToken != nil {
		return "", errToken
	}

	return token, nil
}
