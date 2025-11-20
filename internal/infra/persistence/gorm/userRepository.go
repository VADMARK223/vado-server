package gorm

import (
	"fmt"
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
	return r.db.Transaction(func(tx *gorm.DB) error {
		entity := &UserEntity{Username: u.Username, Password: u.Password, Email: u.Email, Role: u.Role, Color: u.Color}
		if err := tx.Create(entity).Error; err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		u.ID = entity.ID
		return nil
	})
}

func (r *UserRepository) DeleteUser(id uint) error {
	if err := r.db.Delete(&UserEntity{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

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
		Role:      entity.Role,
		Color:     entity.Color,
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
		Role:      entity.Role,
		Color:     entity.Color,
	}, nil
}

func (r *UserRepository) GetAll() ([]user.User, error) {
	var entities []UserEntity
	err := r.db.Find(&entities).Error
	result := make([]user.User, 0, len(entities))

	for _, entity := range entities {
		result = append(result, user.User{
			ID:        entity.ID,
			Username:  entity.Username,
			Email:     entity.Email,
			CreatedAt: entity.CreatedAt,
			Role:      entity.Role,
			Color:     entity.Color,
		})
	}

	return result, err
}
