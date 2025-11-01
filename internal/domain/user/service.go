package user

import (
	"errors"
	"fmt"
	"time"
	"vado_server/internal/domain/auth"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	repo       Repository
	tokenTTL   time.Duration
	refreshTTL time.Duration
}

func NewService(repo Repository, tokenTTL time.Duration, refreshTTL time.Duration) *Service {
	return &Service{
		repo:       repo,
		tokenTTL:   tokenTTL,
		refreshTTL: refreshTTL,
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

func (s *Service) DeleteUser(id uint) error {
	return s.repo.DeleteUser(id)
}

func (s *Service) Login(username, password, secret string) (*User, string, string, error) {
	u, errGetUser := s.repo.GetByUsername(username)
	if errGetUser != nil {
		return nil, "", "", errors.New("user not found")
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) != nil {
		return u, "", "", errors.New("incorrect password")
	}

	accessToken, errToken := auth.CreateToken(u.ID, []string{"user"}, s.tokenTTL, secret)
	if errToken != nil {
		return u, "", "", errors.New(fmt.Sprintf("Error creating token (access): %s", errToken.Error()))
	}

	refreshToken, errToken := auth.CreateToken(u.ID, []string{"user"}, s.refreshTTL, secret)
	if errToken != nil {
		return u, "", "", errors.New(fmt.Sprintf("Error creating token (refresh): %s", errToken.Error()))
	}

	return u, accessToken, refreshToken, nil
}

func (s *Service) Refresh(token string, secret string) (*User, string, error) {
	claims, errParseToken := auth.ParseToken(token, secret)
	if errParseToken != nil {
		return nil, "", status.Error(codes.Unauthenticated, "token read error")
	}
	u, errGetUser := s.repo.GetByID(claims.UserID)
	if errGetUser != nil {
		return nil, "", errors.New("user not found")
	}

	newToken, errToken := auth.CreateToken(u.ID, []string{"user"}, s.tokenTTL, secret)
	if errToken != nil {
		return nil, "", status.Error(codes.Unauthenticated, fmt.Sprintf("Error creating new token: %s", errToken.Error()))
	}

	return u, newToken, nil
}

func (s *Service) GetAllUsersWithRoles() ([]WithRoles, error) {
	return s.repo.GetAllWithRoles()
}
