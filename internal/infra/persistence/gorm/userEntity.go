package gorm

import (
	"time"
)

type UserEntity struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Role      string    `gorm:"not null"`
	Color     string    `gorm:"not null"`
	Email     string    `gorm:"unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (UserEntity) TableName() string {
	return "users"
}
