package user

type Repository interface {
	CreateUser(user User) error
	GetByID(id uint) (*User, error)
	GetByUsername(username string) (*User, error)
	Update(user User) error
	Delete(id uint) error
}
