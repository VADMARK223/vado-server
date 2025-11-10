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

type DTO struct {
	Username string
	Password string
	Email    string
}

type WithRoles struct {
	User
	Roles []RoleDTO
}

type RoleDTO struct {
	ID   uint
	Name string
}
