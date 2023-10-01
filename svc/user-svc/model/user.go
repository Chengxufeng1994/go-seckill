package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"size:255"`
	Password string
	Email    string `gorm:"uniqueIndex"`
	Age      int
}
