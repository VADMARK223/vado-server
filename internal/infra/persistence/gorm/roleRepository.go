package gorm

import (
	"vado_server/internal/app/context"
	"vado_server/internal/domain/role"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoleRepo struct {
	db  *gorm.DB
	log *zap.SugaredLogger
}

func (r RoleRepo) GetAll() ([]role.Role, error) {
	var entities []RoleEntity
	err := r.db.Find(&entities).Error

	result := make([]role.Role, 0, len(entities))
	for _, entity := range entities {
		result = append(result, role.Role{
			ID:   entity.ID,
			Name: entity.Name,
		})
	}

	return result, err
}

func NewRoleRepo(ctx *context.AppContext) role.Repository {
	return &RoleRepo{
		db:  ctx.DB,
		log: ctx.Log,
	}
}
