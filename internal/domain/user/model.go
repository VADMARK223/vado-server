package user

import (
	"time"
)

type User struct {
	ID        uint
	Login     string
	Password  string
	Email     string
	Role      Role
	Color     string
	CreatedAt time.Time

	TasksIDs []uint
}

func New(login, password, email, color string, role Role) User {
	return User{
		Login:    login,
		Password: password,
		Email:    email,
		Role:     role,
		Color:    color,
	}
}

type DTO struct {
	Login    string
	Password string
	Email    string
	Role     Role
	Color    string
}
