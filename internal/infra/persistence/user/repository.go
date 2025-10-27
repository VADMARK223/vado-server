package user

import (
	"vado_server/internal/app/context"
	"vado_server/internal/domain/user"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GormRepository struct {
	db  *gorm.DB
	log *zap.SugaredLogger
}

func NewGormRepo(ctx *context.AppContext) user.Repository {
	return &GormRepository{db: ctx.DB, log: ctx.Log}
}

func (r *GormRepository) CreateUser(u user.User) error {
	entity := &Entity{
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

func (r *GormRepository) GetByID(id uint) (*user.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *GormRepository) GetByUsername(username string) (*user.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *GormRepository) Update(user user.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *GormRepository) Delete(id uint) error {
	//TODO implement me
	panic("implement me")
}
