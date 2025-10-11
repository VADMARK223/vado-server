package repository

import (
	"vado_server/internal/models"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetAll() ([]models.Role, error) {
	var result []models.Role
	err := r.db.Find(&result).Error
	return result, err
}
