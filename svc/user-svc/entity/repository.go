package entity

type Repository interface {
	GetUser(userid uint) (*User, error)
	GetUsers() ([]*User, error)
	GetUserByUsername(username string) (*User, error)
	CreateUser(user *User) (uint, error)
	UpdateUser(user *User) error
	DeleteUser(userid uint) error
}
