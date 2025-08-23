package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	ID            uint64 `gorm:"primaryKey"`
	AuthorID      uint64
	Author        User   `gorm:"foreignKey:AuthorID"`
	Title         string `gorm:"size:100"`
	PlayURL       string `gorm:"size:255"`
	CoverURL      string `gorm:"size:255"`
	FavoriteCount uint32 `gorm:"default:0"`
	CommentCount  uint32 `gorm:"default:0"`
	Duration      uint32
	Resolution    string          `gorm:"size:20"`
	Tags          *datatypes.JSON `gorm:"type:json"`
	Status        int8            `gorm:"default:1"`
}

// TableName 实现接口返回自定义表名
func (Video) TableName() string {
	return "video" // 指定实际表名
}
