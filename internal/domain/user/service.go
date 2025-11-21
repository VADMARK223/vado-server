package user

import (
	"errors"
	"fmt"
	"vado_server/internal/domain/auth"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	repo   Repository
	tokens auth.TokenProvider
}

func NewService(repo Repository, tokens auth.TokenProvider) *Service {
	return &Service{
		repo:   repo,
		tokens: tokens,
	}
}

func (s *Service) CreateUser(dto DTO) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	user := New(dto.Login, string(hash), dto.Email, dto.Username, dto.Color, dto.Role)
	return s.repo.CreateUser(user)
}

func (s *Service) DeleteUser(id uint) error {
	return s.repo.DeleteUser(id)
}

func (s *Service) GetByID(id uint) (User, error) {
	return s.repo.GetByID(id)
}

func (s *Service) Login(login, password string) (User, *auth.TokenPair, error) {
	u, errGetUser := s.repo.GetByLogin(login)

	if errGetUser != nil {
		return User{}, nil, errors.New("user not found")
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		return u, nil, errors.New("incorrect password")
	}

	tokens, err := s.tokens.CreateTokenPair(u.ID, string(u.Role))
	if err != nil {
		return User{}, nil, fmt.Errorf("error creating tokens: %s", err.Error())
	}

	return u, tokens, nil
}

func (s *Service) Refresh(token string) (User, string, error) {
	claims, errParseToken := s.tokens.ParseToken(token)
	if errParseToken != nil {
		return User{}, "", status.Error(codes.Unauthenticated, "token read error")
	}
	u, errGetUser := s.repo.GetByID(claims.UserID())
	if errGetUser != nil {
		return User{}, "", errors.New("user not found")
	}

	newToken, errToken := s.tokens.CreateToken(u.ID, string(u.Role))
	if errToken != nil {
		return User{}, "", status.Error(codes.Unauthenticated, fmt.Sprintf("Error creating new token: %s", errToken.Error()))
	}

	return u, newToken, nil
}

func (s *Service) GetAll() ([]User, error) {
	return s.repo.GetAll()
}
