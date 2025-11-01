package user

type Repository interface {
	CreateUser(user User) error
	DeleteUser(id uint) error
	GetByUsername(username string) (*User, error)
	GetByID(id uint) (*User, error)
	GetAllWithRoles() ([]WithRoles, error)
}
