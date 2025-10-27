package role

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetAll() ([]Role, error) {
	var result []Role
	err := r.db.Find(&result).Error
	return result, err
}
