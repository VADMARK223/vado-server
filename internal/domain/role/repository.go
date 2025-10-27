package role

type Repository interface {
	GetAll() ([]Role, error)
}
