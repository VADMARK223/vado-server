package gorm

import (
	"time"
)

type UserEntity struct {
	ID        uint      `gorm:"primaryKey"`
	Login     string    `gorm:"unique;not null"`
	Username  string    `gorm:"not null"`
	Password  string    `gorm:"not null"`
	Role      string    `gorm:"not null"`
	Color     string    `gorm:"not null"`
	Email     string    `gorm:"unique;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (UserEntity) TableName() string {
	return "users"
}
