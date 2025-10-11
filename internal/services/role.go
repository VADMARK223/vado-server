package services

import (
	"vado_server/internal/models"
	"vado_server/internal/repository"
)

type RoleService struct {
	repo *repository.RoleRepository
}

func NewRoleService(repo *repository.RoleRepository) *RoleService {
	return &RoleService{repo: repo}
}

func (s *RoleService) GetAll() ([]models.Role, error) {
	return s.repo.GetAll()
}
