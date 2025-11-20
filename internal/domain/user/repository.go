package user

type Repository interface {
	CreateUser(user User) error
	DeleteUser(id uint) error
	GetByLogin(login string) (User, error)
	GetByID(id uint) (User, error)
	GetAll() ([]User, error)
}
