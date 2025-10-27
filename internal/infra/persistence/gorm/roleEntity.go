package gorm

type RoleEntity struct {
	ID   int
	Name string `gorm:"unique;not null"`
}

func (RoleEntity) TableName() string {
	return "roles"
}
