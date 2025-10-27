package user

import (
	"time"
)

type User struct {
	ID        uint
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time

	RolesIDs []uint
	TasksIDs []uint
}
