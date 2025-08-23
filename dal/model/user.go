package model

import (
	"gorm.io/datatypes" // 必须引入
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID          uint64 `gorm:"primaryKey"`
	Username    string `gorm:"size:32;uniqueIndex"`
	Password    string `gorm:"size:255"`
	Mobile      string `gorm:"size:20;unique"`
	Email       string `gorm:"size:100;unique"`
	AvatarURL   string `gorm:"size:255"`
	Signature   string `gorm:"size:255"`
	Gender      int8   `gorm:"default:0"`
	Birthday    *time.Time
	IsCreator   bool            `gorm:"default:false"`
	CreatorInfo *datatypes.JSON `gorm:"type:json"`
	Videos      []Video         `gorm:"foreignKey:AuthorID"`
}

// TableName 实现接口返回自定义表名
func (User) TableName() string {
	return "douyin_user" // 指定实际表名
}
