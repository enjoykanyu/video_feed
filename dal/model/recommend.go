package model

import (
	"time"
)

// UserInterest 用户兴趣标签表
type UserInterest struct {
	ID        int64   `gorm:"primaryKey"`
	UserID    int64   `gorm:"index"`
	TagName   string  `gorm:"size:50;index"`
	Weight    float64 `gorm:"default:0"`
	UpdatedAt time.Time
}

// WatchHistory 用户观看历史表
type WatchHistory struct {
	ID            int64 `gorm:"primaryKey"`
	UserID        int64 `gorm:"index"`
	VideoID       int64 `gorm:"index"`
	WatchDuration int32 `gorm:"comment:观看时长(秒)"`
	CreatedAt     time.Time
}

// VideoTag 视频标签表
type VideoTag struct {
	ID      int64   `gorm:"primaryKey"`
	VideoID int64   `gorm:"index"`
	TagName string  `gorm:"size:50;index"`
	Weight  float64 `gorm:"default:0"`
}

// UserProfile 用户画像表
type UserProfile struct {
	ID            int64 `gorm:"primaryKey"`
	UserID        int64 `gorm:"uniqueIndex"`
	LastActiveAt  time.Time
	TotalWatch    int64 `gorm:"default:0;comment:总观看时长"`
	TotalLikes    int64 `gorm:"default:0;comment:总点赞数"`
	TotalComments int64 `gorm:"default:0;comment:总评论数"`
	TotalShares   int64 `gorm:"default:0;comment:总分享数"`
	UpdatedAt     time.Time
}

// VideoScore 视频评分表
type VideoScore struct {
	ID           int64   `gorm:"primaryKey"`
	VideoID      int64   `gorm:"uniqueIndex"`
	HotScore     float64 `gorm:"default:0;comment:热度分数"`
	QualityScore float64 `gorm:"default:0;comment:质量分数"`
	UpdatedAt    time.Time
}
