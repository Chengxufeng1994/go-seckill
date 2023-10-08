package model

import (
	"github.com/Chengxufeng1994/go-seckill/svc/user-svc/entity"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Password string
	Email    string
	Age      int
}

func UserDao2Dto(dto *entity.User) *User {

	var dao *User
	dao.ID = dto.ID
	dao.CreatedAt = dto.CreatedAt
	dao.UpdatedAt = dto.UpdatedAt
	dao.DeletedAt = dto.DeletedAt
	dao.Username = dto.Username
	dao.Password = dto.Password
	dao.Email = dto.Password
	dao.Age = dto.Age

	return dao
}
