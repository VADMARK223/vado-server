package task

type Repository interface {
	GetAll() ([]Task, error)
	GetAllByUserID(ID uint) ([]Task, error)
}
