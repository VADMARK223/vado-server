package gorm

import (
	"vado_server/internal/app/context"
	"vado_server/internal/domain/user"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db  *gorm.DB
	log *zap.SugaredLogger
}

func NewUserRepo(ctx *context.AppContext) user.Repository {
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
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) Update(user user.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) Delete(id uint) error {
	//TODO implement me
	panic("implement me")
}
