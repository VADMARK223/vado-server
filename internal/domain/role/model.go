package role

type Role struct {
	ID   uint
	Name string
}

var (
	Admin     = Role{ID: 1, Name: "admin"}
	Moderator = Role{ID: 2, Name: "moderator"}
	User      = Role{ID: 3, Name: "user"}
)

var All = []Role{Admin, Moderator, User}
