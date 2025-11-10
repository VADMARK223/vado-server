package gorm

type RoleEntity struct {
	ID   uint
	Name string `gorm:"unique;not null"`
}

func (RoleEntity) TableName() string {
	return "roles"
}
