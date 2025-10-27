package user

import (
	"time"
	"vado_server/internal/domain/role"
	"vado_server/internal/domain/task"
)

type Entity struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Email     string    `gorm:"unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Связь с задачами (1 пользователь и много задач)
	Tasks []task.Task `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`

	Roles []role.Role `gorm:"many2many:user_roles;constraint:OnDelete:CASCADE;"`
}

func (Entity) TableName() string {
	return "users"
}
