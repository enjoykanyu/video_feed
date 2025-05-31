package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int64  `gorm:"primaryKey"` // 替代uint
	Phone    string `gorm:"unique"`
	Password string
	Username string
	Avatar   string
}

// TableName 实现接口返回自定义表名
func (User) TableName() string {
	return "user_db" // 指定实际表名
}
