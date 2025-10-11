package models

type Role struct {
	ID   int
	Name string
}

func (r Role) TableName() string {
	return "roles"
}
