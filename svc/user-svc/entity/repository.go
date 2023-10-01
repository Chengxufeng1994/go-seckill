package entity

import "github.com/Chengxufeng1994/go-seckill/svc/user-svc/model"

type Repository interface {
	GetUser(userid uint) (model.User, error)
	GetUsers() ([]model.User, error)
	GetUserByUsername(username string) (model.User, error)
	CreateUser(user model.User) (uint, error)
	UpdateUser(user model.User) error
	DeleteUser(userid uint) error
}
