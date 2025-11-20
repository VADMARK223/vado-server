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
	entity := toEntity(u)
	if err := r.db.Create(&entity).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) DeleteUser(id uint) error {
	if err := r.db.Delete(&UserEntity{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByUsername(username string) (user.User, error) {
	var entity UserEntity
	if err := r.db.Where("username = ?", username).First(&entity).Error; err != nil {
		return user.User{}, err
	}

	return toDomain(entity), nil
}

func (r *UserRepository) GetByID(id uint) (user.User, error) {
	var entity UserEntity
	if err := r.db.First(&entity, id).First(&entity).Error; err != nil {
		return user.User{}, err
	}

	return toDomain(entity), nil
}

func (r *UserRepository) GetAll() ([]user.User, error) {
	var entities []UserEntity
	err := r.db.Find(&entities).Error
	result := make([]user.User, 0, len(entities))

	for _, entity := range entities {
		result = append(result, toDomain(entity))
	}

	return result, err
}

func toDomain(e UserEntity) user.User {
	return user.User{
		ID:        e.ID,
		Login:     e.Username,
		Password:  e.Password,
		Email:     e.Email,
		CreatedAt: e.CreatedAt,
		Role:      user.Role(e.Role),
		Color:     e.Color,
	}
}

func toEntity(u user.User) UserEntity {
	return UserEntity{
		ID:        u.ID,
		Username:  u.Login,
		Password:  u.Password,
		Email:     u.Email,
		Color:     u.Color,
		Role:      string(u.Role), // доменный тип → строка
		CreatedAt: u.CreatedAt,
	}
}
