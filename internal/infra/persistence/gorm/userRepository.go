package gorm

import (
	"vado_server/internal/app"
	"vado_server/internal/domain/user"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db  *gorm.DB
	log *zap.SugaredLogger
}

func NewUserRepo(ctx *app.Context) user.Repository {
	return &UserRepository{
		db:  ctx.DB,
		log: ctx.Log,
	}
}

func (r *UserRepository) CreateUser(u user.User) error {
	entity := &UserEntity{
		Username: u.Username,
		Password: u.Password,
		Email:    u.Email,
	}

	if err := r.db.Create(entity).Error; err != nil {
		r.log.Errorw("failed to create user", "error", err)
		return err
	}

	u.ID = entity.ID
	return nil
}

func (r *UserRepository) GetByUsername(username string) (*user.User, error) {
	var entity UserEntity
	if err := r.db.Where("username = ?", username).First(&entity).Error; err != nil {
		return nil, err
	}

	return &user.User{
		ID:        entity.ID,
		Username:  entity.Username,
		Password:  entity.Password,
		Email:     entity.Email,
		CreatedAt: entity.CreatedAt,
	}, nil
}

func (r *UserRepository) GetByID(id uint) (*user.User, error) {
	var entity UserEntity
	if err := r.db.First(&entity, id).First(&entity).Error; err != nil {
		return nil, err
	}

	return &user.User{
		ID:        entity.ID,
		Username:  entity.Username,
		Password:  entity.Password,
		Email:     entity.Email,
		CreatedAt: entity.CreatedAt,
	}, nil
}

func (r *UserRepository) GetAllWithRoles() ([]user.WithRoles, error) {
	var entities []UserEntity

	if err := r.db.Preload("Roles").Find(&entities).Error; err != nil {
		r.log.Errorw("failed to get users with roles", "error", err)
		return nil, err
	}

	result := make([]user.WithRoles, 0, len(entities))

	for _, entity := range entities {
		roles := make([]user.RoleDTO, 0, len(entity.Roles))
		for _, roleEntity := range entity.Roles {
			roles = append(roles, user.RoleDTO{
				ID:   roleEntity.ID,
				Name: roleEntity.Name,
			})
		}

		result = append(result, user.WithRoles{
			User: user.User{
				ID:        entity.ID,
				Username:  entity.Username,
				Email:     entity.Email,
				CreatedAt: entity.CreatedAt,
			},
			Roles: roles,
		})
	}

	return result, nil
}
