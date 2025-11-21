package user

import (
	"time"
)

type User struct {
	ID        uint
	Login     string
	Username  string
	Password  string
	Email     string
	Role      Role
	Color     string
	CreatedAt time.Time

	TasksIDs []uint
}

func (u User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func New(login, password, email, username, color string, role Role) User {
	return User{
		Login:    login,
		Password: password,
		Email:    email,
		Username: username,
		Role:     role,
		Color:    color,
	}
}

type DTO struct {
	Login    string
	Password string
	Email    string
	Username string
	Role     Role
	Color    string
}
