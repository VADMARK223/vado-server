package services

import (
	"vado_server/internal/models"
	"vado_server/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(dto models.UserDTO) error {
	return s.repo.CreateUser(dto)
}
