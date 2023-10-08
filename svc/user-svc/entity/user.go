package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"size:255;unique"`
	Password string `gorm:"size:255;unique"`
	Email    string `gorm:"size:255;unique"`
	Age      int    `gorm:"size:10"`
}
