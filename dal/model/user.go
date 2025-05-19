package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Phone    string `gorm:"unique"`
	Password string
	Nickname string
	Avatar   string
}
