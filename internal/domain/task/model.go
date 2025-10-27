package task

import (
	"time"
)

type Task struct {
	ID          uint
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Completed   bool
}
