package task

import (
	"time"
)

type Task struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Completed   bool      `gorm:"default:false"`

	UserID uint `gorm:"not null"` // FK
}
