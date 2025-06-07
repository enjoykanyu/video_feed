package model

import (
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	ID          int64 `gorm:"primaryKey"` // 替代uint
	UserID      int64
	Title       string
	Description string
	CoverURL    string
	VideoURL    string
	Status      int32
}

// TableName 实现接口返回自定义表名
func (Video) TableName() string {
	return "video" // 指定实际表名
}
