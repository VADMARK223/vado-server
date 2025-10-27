package gorm

import "time"

type TaskEntity struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Completed   bool      `gorm:"default:false"`

	UserID uint
}

func (TaskEntity) TableName() string {
	return "tasks"
}
