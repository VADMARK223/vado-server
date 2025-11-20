package user

import (
	"time"
)

type User struct {
	ID        uint
	Username  string
	Password  string
	Email     string
	Role      Role
	Color     string
	CreatedAt time.Time

	TasksIDs []uint
}

type DTO struct {
	Username string
	Password string
	Email    string
	Role     Role
	Color    string
}
