package user

import (
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
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
