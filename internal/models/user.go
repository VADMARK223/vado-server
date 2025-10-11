package models

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	Email     string    `gorm:"unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Связь с задачами (1 пользователь и много задач)
	Tasks []Task `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

func (u User) TableName() string {
	return "users"
}
