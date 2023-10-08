package model

import (
	"database/sql"
	"github.com/Chengxufeng1994/go-seckill/svc/user-svc/entity"
	"time"
)

type User struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
	Username  string
	Password  string
	Email     string
	Age       int
}

func UserDao2Dto(dto *entity.User) *User {
	return &User{
		ID:        dto.ID,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		Username:  dto.Username,
		Password:  dto.Password,
		Email:     dto.Password,
		Age:       dto.Age,
	}
}
