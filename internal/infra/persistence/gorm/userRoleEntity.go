package gorm

type UserRoleEntity struct {
	UserID uint `gorm:"not null"`
	RoleID uint `gorm:"not null"`
}

func (u UserRoleEntity) TableName() string {
	return "user_roles"
}
