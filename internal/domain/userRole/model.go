package userRole

type UserRole struct {
	UserID uint `gorm:"not null"`
	RoleID int  `gorm:"not null"`
}

func (u UserRole) TableName() string {
	return "user_roles"
}
