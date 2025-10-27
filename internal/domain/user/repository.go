package user

import (
	"vado_server/internal/config/role"
	"vado_server/internal/domain/userRole"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(dto DTO) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	user := User{
		Username: dto.Username,
		Password: string(hash),
		Email:    dto.Email,
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		ur := userRole.UserRole{
			UserID: user.ID,
			RoleID: role.User,
		}

		if err := tx.Create(&ur).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}
