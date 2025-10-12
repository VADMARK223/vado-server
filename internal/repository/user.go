package repository

import (
	"vado_server/internal/constants/role"
	"vado_server/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(dto models.UserDTO) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	user := models.User{
		Username: dto.Username,
		Password: string(hash),
		Email:    dto.Email,
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		userRole := models.UserRole{
			UserID: user.ID,
			RoleID: role.User,
		}

		if err := tx.Create(&userRole).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}
